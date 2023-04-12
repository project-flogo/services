package postgres

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/project-flogo/core/data/coerce"
	metadata2 "github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/store/metadata"
	"github.com/project-flogo/services/flow-state/store/task"
)

func NewStore(settings map[string]interface{}) (*StepStore, error) {
	db, err := NewDB(settings)
	dbDetails := &DBDetails{}
	if err != nil {
		dbDetails.Connected = false
		dbDetails.Message = err.Error()
		dbDetails.SmVersion = "1.0"
		return &StepStore{db: &StatefulDB{db: db, dbDetails: dbDetails}, settings: settings}, nil
	}

	statefulDB := &StatefulDB{db: db}
	dbDetails.getDBDetails(statefulDB)
	statefulDB.dbDetails = dbDetails
	return &StepStore{db: statefulDB, settings: settings}, err

}

type DBDetails struct {
	SmVersion    string `json:"smVersion"`
	Connected    bool   `json:"connected"`
	TablesExists bool   `json:"tablesExists"`
	Message      string `json:"message"`
	Status       bool   `json:"status"`
}

type StepStore struct {
	sync.RWMutex
	db             *StatefulDB
	userId         string
	appId          string
	stepContainers map[string]*stepContainer
	snapshots      sync.Map
	settings       map[string]interface{}
}

func (d *DBDetails) getDBDetails(db *StatefulDB) {
	sqlSelect := "SELECT count(*) FROM information_schema.tables WHERE table_name in ('flowstate' , 'appstate' , 'steps')"
	set, err := db.query(sqlSelect, nil)
	if err != nil {
		d.TablesExists = false
		d.Message = "One or more required tables are missing in the database schema"
	} else {
		for _, v := range set.Record {
			m := *v
			count, _ := coerce.ToInt(m["count"])
			if count == 3 {
				d.TablesExists = true
			} else {
				d.TablesExists = false
				d.Message = "One or more required tables are missing in the database schema"

			}
		}
	}
	if !d.TablesExists {
		return
	}
	sqlSelect = "SELECT column_name FROM information_schema.columns WHERE table_name = 'flowstate' and column_name = 'flowinput'"
	set, err = db.query(sqlSelect, nil)
	if err != nil {
		d.SmVersion = "1.0"
	} else {
		if set.Record == nil {
			d.SmVersion = "1.0"
		}
		if len(set.Record) > 0 {
			d.SmVersion = "2.0"
		}
	}
}

func (s *StepStore) Status() interface{} {

	if !s.db.dbDetails.Connected {
		return s.db.dbDetails
	}

	if err := s.db.db.Ping(); err != nil {
		s.db.dbDetails.Status = false
	} else {
		s.db.dbDetails.Status = true
	}

	return s.db.dbDetails
}

func (s *StepStore) MaxConcurrencyLimit() int {
	if s.db.dbDetails.Connected {
		return s.db.db.Stats().MaxOpenConnections
	} else {
		return 1
	}

}

func (s *StepStore) GetStatus(flowId string) int {
	s.RLock()
	sc, ok := s.stepContainers[flowId]
	s.RUnlock()

	if ok {
		return sc.Status()
	}
	return -1
}

//func (s *StepStore) GetFlow(flowId string) *state.FlowInfo {
//
//	s.RLock()
//	sc, ok := s.stepContainers[flowId]
//	s.RUnlock()
//
//	if ok {
//		return &state.FlowInfo{Id: flowId, Status: sc.Status(), FlowURI: sc.flowURI}
//	}
//
//	return nil
//}

func (s *StepStore) GetFailedFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error) {
	var whereStr = "where"
	if len(metadata.Username) < 0 {
		return nil, fmt.Errorf("unauthorized, please provide user infomation")
	}
	whereStr += " userId='" + metadata.Username + "'"

	if len(metadata.AppName) < 0 {
		return nil, fmt.Errorf("please provide App name")
	}
	whereStr += "  and appName='" + metadata.AppName + "'"

	if len(metadata.AppVersion) < 0 {
		return nil, fmt.Errorf("please provide App Version")
	}
	whereStr += "  and appVersion='" + metadata.AppVersion + "'"

	if len(metadata.HostId) > 0 {
		whereStr += "  and hostId='" + metadata.HostId + "'"
	}

	if len(metadata.FlowName) > 0 {
		whereStr += "  and flowname='" + metadata.FlowName + "'"
	}

	whereStr += " and status = 'Failed'"

	set, err := s.db.query("select flowinstanceid, flowname, status from flowstate "+whereStr, nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetFailedFlows after successful connection retry  ")
				set, err = s.db.query("select flowinstanceid, flowname, status from flowstate "+whereStr, nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return nil, err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return nil, retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return nil, err
		}
	}

	var flowinfo []*state.FlowInfo
	for _, v := range set.Record {
		m := *v
		id, _ := coerce.ToString(m["flowinstanceid"])
		flowName, _ := coerce.ToString(m["flowname"])
		status, _ := coerce.ToString(m["status"])
		info := &state.FlowInfo{
			Id:         id,
			FlowName:   flowName,
			FlowStatus: status,
		}
		flowinfo = append(flowinfo, info)
	}
	return flowinfo, err
}

func (s *StepStore) GetCompletedFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error) {
	var whereStr = "where"
	if len(metadata.Username) < 0 {
		return nil, fmt.Errorf("unauthorized, please provide user infomation")
	}
	whereStr += " userId='" + metadata.Username + "'"

	if len(metadata.AppName) < 0 {
		return nil, fmt.Errorf("please provide App name")
	}
	whereStr += "  and appName='" + metadata.AppName + "'"

	if len(metadata.AppVersion) < 0 {
		return nil, fmt.Errorf("please provide App Version")
	}
	whereStr += "  and appVersion='" + metadata.AppVersion + "'"

	if len(metadata.HostId) > 0 {
		whereStr += "  and hostId='" + metadata.HostId + "'"
	}

	if len(metadata.FlowName) > 0 {
		whereStr += "  and flowname='" + metadata.FlowName + "'"
	}

	whereStr += " and status = 'Completed'"

	set, err := s.db.query("select flowinstanceid, flowname, status from flowstate "+whereStr, nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetCompletedFlows after successful connection retry  ")
				set, err = s.db.query("select flowinstanceid, flowname, status from flowstate "+whereStr, nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return nil, err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return nil, retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return nil, err
		}
	}

	var flowinfo []*state.FlowInfo
	for _, v := range set.Record {
		m := *v
		id, _ := coerce.ToString(m["flowinstanceid"])
		flowName, _ := coerce.ToString(m["flowname"])
		status, _ := coerce.ToString(m["status"])
		info := &state.FlowInfo{
			Id:         id,
			FlowName:   flowName,
			FlowStatus: status,
		}
		flowinfo = append(flowinfo, info)
	}
	return flowinfo, err
}

func (s *StepStore) GetFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error) {
	var whereStr = "where"
	if len(metadata.Username) < 0 {
		return nil, fmt.Errorf("unauthorized, please provide user infomation")
	}
	whereStr += " userId='" + metadata.Username + "'"

	if len(metadata.AppName) < 0 {
		return nil, fmt.Errorf("please provide App name")
	}
	whereStr += "  and appName='" + metadata.AppName + "'"

	if len(metadata.AppVersion) < 0 {
		return nil, fmt.Errorf("please provide App Version")
	}
	whereStr += "  and appVersion='" + metadata.AppVersion + "'"

	if len(metadata.HostId) > 0 {
		whereStr += "  and hostId='" + metadata.HostId + "'"
	}

	if len(metadata.FlowName) > 0 {
		whereStr += "  and flowname='" + metadata.FlowName + "'"
	}

	if len(metadata.Status) > 0 {
		whereStr += "  and status='" + metadata.Status + "'"
	}

	if len(metadata.Offset) > 0 && len(metadata.Limit) > 0 {
		offsetLimitStr := "  offset '" + metadata.Offset + "'  limit  '" + metadata.Limit + "'"
		whereStr += offsetLimitStr
	}

	set, err := s.db.query("select flowinstanceid, flowname, status, hostid, starttime, endtime from flowstate "+whereStr, nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetFlows after successful connection retry  ")
				set, err = s.db.query("select flowinstanceid, flowname, status, hostid, starttime, endtime from flowstate "+whereStr, nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return nil, err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return nil, retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return nil, err
		}
	}

	var flowinfo []*state.FlowInfo
	for _, v := range set.Record {
		m := *v
		id, _ := coerce.ToString(m["flowinstanceid"])
		flowName, _ := coerce.ToString(m["flowname"])
		status, _ := coerce.ToString(m["status"])
		hostid, _ := coerce.ToString(m["hostid"])
		starttime, _ := coerce.ToString(m["starttime"])
		endtime, _ := coerce.ToString(m["endtime"])
		info := &state.FlowInfo{
			Id:         id,
			FlowName:   flowName,
			HostId:     hostid,
			FlowStatus: status,
			StartTime:  starttime,
			EndTime:    endtime,
		}
		flowinfo = append(flowinfo, info)
	}
	return flowinfo, err
}

func (s *StepStore) GetFlowsWithRecordCount(mtdata *metadata.Metadata) (*metadata.FlowRecord, error) {
	var whereStr = "where"
	if len(mtdata.Username) < 0 {
		return nil, fmt.Errorf("unauthorized, please provide user infomation")
	}
	whereStr += " userId='" + mtdata.Username + "'"

	if len(mtdata.AppName) < 0 {
		return nil, fmt.Errorf("please provide App name")
	}
	whereStr += "  and appName='" + mtdata.AppName + "'"

	if len(mtdata.AppVersion) < 0 {
		return nil, fmt.Errorf("please provide App Version")
	}
	whereStr += "  and appVersion='" + mtdata.AppVersion + "'"

	if len(mtdata.HostId) > 0 {
		whereStr += "  and hostId='" + mtdata.HostId + "'"
	}

	if len(mtdata.FlowName) > 0 {
		whereStr += "  and flowname='" + mtdata.FlowName + "'"
	}

	if len(mtdata.Status) > 0 {
		whereStr += "  and status='" + mtdata.Status + "'"
	}

	if len(mtdata.FlowInstanceId) > 0 {
		whereStr += "  and (flowinstanceid='" + mtdata.FlowInstanceId + "' or rerunofflowinstanceid='" + mtdata.FlowInstanceId + "' )"
	}

	if len(mtdata.Interval) > 0 {
		whereStr += "  and starttime >= NOW() - INTERVAL '" + mtdata.Interval + "'"
	}

	if len(mtdata.StartTime) > 0 && len(mtdata.EndTime) > 0 {
		whereStr += "  and starttime >= '" + mtdata.StartTime + "' and starttime <= '" + mtdata.EndTime + "'"
	}

	whereStr += " order by starttime desc"

	if len(mtdata.Offset) > 0 && len(mtdata.Limit) > 0 {
		offsetLimitStr := "  offset '" + mtdata.Offset + "'  limit  '" + mtdata.Limit + "'"
		whereStr += offsetLimitStr
	}

	var set *ResultSet
	var err error
	if s.db.dbDetails.SmVersion == "1.0" {
		set, err = s.db.query("select flowinstanceid, flowname, status, hostid, starttime, endtime, executiontime, rerunofflowinstanceid, count(*) over() AS full_count from flowstate "+whereStr, nil)
	} else {
		set, err = s.db.query("select flowinstanceid, flowname, status, hostid, starttime, endtime, executiontime, rerunofflowinstanceid, reruncount, flowinput, count(*) over() AS full_count from flowstate "+whereStr, nil)
	}

	if err != nil {
		pqerror, ok := err.(*pq.Error)
		if ok {
			if pqerror.Routine == "errorMissingColumn" {
				set, err = s.db.query("select flowinstanceid, flowname, status, hostid, starttime, endtime, executiontime, rerunofflowinstanceid, count(*) over() AS full_count from flowstate "+whereStr, nil)
			}
		}

		if err != nil {

			if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
				strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
				strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
				strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
				if retryErr := s.RetryDBConnection(); retryErr == nil {
					logCache.Debugf("Retrying from GetFlowsWithRecordCount after successful connection retry  ")
					set, err = s.db.query("select flowinstanceid, flowname, status, hostid, starttime, endtime, executiontime, rerunofflowinstanceid, reruncount, flowinput, count(*) over() AS full_count from flowstate "+whereStr, nil)
					if err != nil {
						logCache.Errorf("Could not connect to database server error:, %s", err.Error())
						return nil, err
					}
				} else {
					logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
					return nil, retryErr
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", err.Error())
				return nil, err
			}
		}
	}
	var count int32
	var flowinfo []*state.FlowInfo
	for _, v := range set.Record {
		m := *v
		id, _ := coerce.ToString(m["flowinstanceid"])
		flowName, _ := coerce.ToString(m["flowname"])
		status, _ := coerce.ToString(m["status"])
		hostid, _ := coerce.ToString(m["hostid"])
		starttime, _ := coerce.ToString(m["starttime"])
		endtime, _ := coerce.ToString(m["endtime"])
		executiontime, _ := coerce.ToString(m["executiontime"])
		count, _ = coerce.ToInt32(m["full_count"])
		originalInstanceId, _ := coerce.ToString(m["rerunofflowinstanceid"])
		reRunCount, _ := coerce.ToInt(m["reruncount"])
		var flowInput map[string]interface{}

		if m["flowinput"] != nil {
			flowInput = make(map[string]interface{})
		}

		info := &state.FlowInfo{
			Id:                 id,
			FlowName:           flowName,
			HostId:             hostid,
			FlowStatus:         status,
			StartTime:          starttime,
			EndTime:            endtime,
			ExecutionTime:      executiontime,
			OriginalInstanceId: originalInstanceId,
			RerunCount:         reRunCount,
			FlowInputs:         flowInput,
		}
		flowinfo = append(flowinfo, info)
	}

	val := &metadata.FlowRecord{
		Count:    count,
		FlowData: flowinfo}

	return val, nil
}

func (s *StepStore) GetFlow(flowid string, metadata *metadata.Metadata) (*state.FlowInfo, error) {
	var whereStr = "where flowinstanceid = '" + flowid + "'"
	if len(metadata.Username) > 0 {
		whereStr += " and userId='" + metadata.Username + "'"
	}

	if len(metadata.AppName) > 0 {
		whereStr += "  and appName='" + metadata.AppName + "'"
	}
	if len(metadata.AppVersion) > 0 {
		whereStr += "  and appVersion='" + metadata.AppVersion + "'"
	}
	if len(metadata.HostId) > 0 {
		whereStr += "  and hostId='" + metadata.HostId + "'"
	}

	set, err := s.db.query("select flowinstanceid, flowname, status, flowinput, reruncount, rerunofflowinstanceid from flowstate "+whereStr, nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetFlow after successful connection retry  ")
				set, err = s.db.query("select flowinstanceid, flowname, status, flowinput from flowstate "+whereStr, nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return nil, err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return nil, retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return nil, err
		}
	}

	var flowinfo []*state.FlowInfo
	for _, v := range set.Record {
		m := *v
		id, _ := coerce.ToString(m["flowinstanceid"])
		flowName, _ := coerce.ToString(m["flowname"])
		status, _ := coerce.ToString(m["status"])

		flowInputBytes, err := coerce.ToBytes(m["flowinput"])
		if err != nil {
			logCache.Errorf("decodeBase64 for flowInputBytes in GetFlow error:, %s", err.Error())
			return nil, fmt.Errorf("decodeBase64 for flowInputBytes in GetFlow error:", err.Error())
		}
		dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(flowInputBytes)))
		n, err := base64.StdEncoding.Decode(dbuf, flowInputBytes)
		if err != nil {
			return nil, err
		}
		stepData := dbuf[:n]
		var flowInput map[string]interface{}
		err = json.Unmarshal(stepData, &flowInput)
		if err != nil {
			return nil, err
		}
		info := &state.FlowInfo{
			Id:         id,
			FlowName:   flowName,
			FlowStatus: status,
			FlowURI:    "res://flow:" + flowName,
			FlowInputs: flowInput,
		}
		flowinfo = append(flowinfo, info)
	}
	if len(flowinfo) <= 0 {
		return nil, fmt.Errorf("flow details [%s] not found", flowid)
	}
	return flowinfo[0], err
}

func (s *StepStore) GetFlowNames(metadata *metadata.Metadata) ([]string, error) {
	var whereStr = "where "
	if len(metadata.Username) > 0 {
		whereStr += " userId='" + metadata.Username + "'"
	}
	if len(metadata.AppName) > 0 {
		whereStr += "  and appName='" + metadata.AppName + "'"
	}
	if len(metadata.AppVersion) > 0 {
		whereStr += "  and appVersion='" + metadata.AppVersion + "'"
	}

	if len(metadata.HostId) > 0 {
		whereStr += "  and hostId='" + metadata.HostId + "'"
	}

	set, err := s.db.query("select distinct(flowname) from flowstate "+whereStr, nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetFlowNames after successful connection retry  ")
				set, err = s.db.query("select distinct(flowname) from flowstate "+whereStr, nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return nil, err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return nil, retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return nil, err
		}
	}

	var flownameArray []string
	for _, v := range set.Record {
		m := *v
		flowName, _ := coerce.ToString(m["flowname"])
		flownameArray = append(flownameArray, flowName)
	}
	return flownameArray, err
}

func (s *StepStore) GetAppVersions(metadata *metadata.Metadata) ([]string, error) {
	var whereStr = "where "
	if len(metadata.Username) > 0 {
		whereStr += " userId='" + metadata.Username + "'"
	}
	if len(metadata.AppName) > 0 {
		whereStr += "  and appname='" + metadata.AppName + "'"
	}

	set, err := s.db.query("select distinct(appVersion) from flowstate "+whereStr, nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetAppVersions after successful connection retry  ")
				set, err = s.db.query("select distinct(appVersion) from flowstate "+whereStr, nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return nil, err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return nil, retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return nil, err
		}
	}

	var appVersionArray []string
	for _, v := range set.Record {
		m := *v
		appVersion, _ := coerce.ToString(m["appversion"])
		appVersionArray = append(appVersionArray, appVersion)
	}
	return appVersionArray, err
}

func (s *StepStore) GetAppState(metadata *metadata.Metadata) (string, error) {
	var whereStr = "where "
	if len(metadata.Username) > 0 {
		whereStr += " userId='" + metadata.Username + "'"
	}
	if len(metadata.AppName) > 0 {
		whereStr += "  and appname='" + metadata.AppName + "'"
	}

	set, err := s.db.query("select persistenceEnabled from appstate "+whereStr, nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetAppState after successful connection retry  ")
				set, err = s.db.query("select persistenceEnabled from appstate  "+whereStr, nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return "", err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return "", retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return "", err
		}
	}
	persistenceEnabled := ""
	for _, v := range set.Record {
		m := *v
		persistenceEnabled, _ = coerce.ToString(m["persistenceenabled"])
	}
	return persistenceEnabled, err
}

func (s *StepStore) SaveAppState(metadata *metadata.Metadata) error {
	_, err := s.db.InsertAppState(metadata)
	if err != nil && (err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
		strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
		strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
		strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout")) {
		if retryErr := s.RetryDBConnection(); retryErr == nil {
			logCache.Debug("Retrying from SaveAppState after successful connection retry  ")
			_, err = s.db.InsertAppState(metadata)
			if err != nil {
				logCache.Errorf("Could not connect to database server error:, %s", err.Error())
				return err
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
			return retryErr
		}
	}
	return err
}

func (s *StepStore) SaveStep(step *state.Step) error {
	_, err := s.db.InsertSteps(step)
	if err != nil && (err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
		strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
		strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
		strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout")) {
		if retryErr := s.RetryDBConnection(); retryErr == nil {
			logCache.Debugf("Retrying from SaveStep after successful connection retry  ")
			_, err = s.db.InsertSteps(step)
			if err != nil {
				logCache.Errorf("Could not connect to database server error:, %s", err.Error())
				return err
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
			return retryErr
		}
	}
	return err
}

func (s *StepStore) DeleteSteps(flowId string, stepId string) error {
	_, err := s.db.DeleteSteps(flowId, stepId)
	if err != nil && (err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
		strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
		strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
		strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout")) {
		if retryErr := s.RetryDBConnection(); retryErr == nil {
			logCache.Debugf("Retrying from  DeleteSteps after successful connection retry  ")
			_, err = s.db.DeleteSteps(flowId, stepId)
			if err != nil {
				logCache.Errorf("Could not connect to database server error:, %s", err.Error())
				return err
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
			return retryErr
		}
	}
	return err
}

func (s *StepStore) GetSteps(flowId string) ([]*state.Step, error) {

	set, err := s.db.query("select stepdata from steps where flowinstanceid = '"+flowId+"'", nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetSteps after successful connection retry  ")
				set, err = s.db.query("select stepdata from steps where flowinstanceid = '"+flowId+"'", nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return nil, err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return nil, retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return nil, err
		}
	}

	var steps []*state.Step
	for _, v := range set.Record {
		m := *v
		s1, err := coerce.ToBytes(m["stepdata"])
		if err != nil {
			logCache.Errorf("decodeBase64 for step data error:, %s", err.Error())
			return nil, fmt.Errorf("decodeBase64 for step data error:", err.Error())
		}
		dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(s1)))
		n, err := base64.StdEncoding.Decode(dbuf, s1)
		stePdata := dbuf[:n]
		var step *state.Step
		err = json.Unmarshal(stePdata, &step)
		if err != nil {
			logCache.Errorf("Marshalling error:, %s", err.Error())
			return nil, err
		}
		steps = append(steps, step)
	}

	if len(steps) <= 0 {
		return nil, fmt.Errorf("step for flow instance [%s] not found", flowId)
	}
	return steps, err
}

func (s *StepStore) GetStepsAsTasks(flowId string) ([][]*task.Task, error) {

	set, err := s.db.query("select stepdata from steps where flowinstanceid = '"+flowId+"'", nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetStepsAsTasks after successful connection retry  ")
				set, err = s.db.query("select stepdata from steps where flowinstanceid = '"+flowId+"'", nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return nil, err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return nil, retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return nil, err
		}
	}

	var steps []*state.Step
	for _, v := range set.Record {
		m := *v
		s1, err := coerce.ToBytes(m["stepdata"])
		if err != nil {
			return nil, fmt.Errorf("decodeBase64 for step data error:", err.Error())
		}
		dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(s1)))
		n, err := base64.StdEncoding.Decode(dbuf, s1)
		stePdata := dbuf[:n]
		var step *state.Step
		err = json.Unmarshal(stePdata, &step)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)

	}

	if len(steps) <= 0 {
		return nil, fmt.Errorf("step for flow instance [%s] not found", flowId)
	}
	var taskValueArray [][]*task.Task
	for _, stepval := range steps {
		taskValue, err := task.StepToTask(stepval)
		if err != nil {
			return nil, err
		}
		taskValueArray = append(taskValueArray, taskValue)
	}

	return taskValueArray, err
}

func (s *StepStore) GetStepdataForActivity(flowId, stepid, taskname string) ([]*task.Task, error) {
	query := "select stepdata from steps where flowinstanceid = '" + flowId + "' and stepid = '" + stepid + "'"
	if taskname != "" {
		query += " and taskname = '" + taskname + "'"
	}
	set, err := s.db.query(query, nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetStepdataForActivity after successful connection retry  ")
				set, err = s.db.query(query, nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return nil, err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return nil, retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return nil, err
		}
	}
	var step *state.Step
	for _, v := range set.Record {
		m := *v
		s1, err := coerce.ToBytes(m["stepdata"])
		if err != nil {
			return nil, fmt.Errorf("decodeBase64 for step data error:", err.Error())
		}
		dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(s1)))
		n, err := base64.StdEncoding.Decode(dbuf, s1)
		if err != nil {
			return nil, err
		}
		stepData := dbuf[:n]
		err = json.Unmarshal(stepData, &step)
		if err != nil {
			return nil, err
		}
	}
	if step == nil {
		return nil, fmt.Errorf("No step data found for matching input")
	}
	taskValue, err := task.StepToTask(step)
	if err != nil {
		return nil, err
	}
	// return taskValue, err
	// if array having 2 tasks + status waiting, and name and subflowid same ( to identify its 'callsubflow' activity), then query stepdata for its enclosing record (failed or completed)
	//if present and merge the output
	if len(taskValue) == 2 {
		firsttask := taskValue[0]
		taskname1 := firsttask.Id
		subflowid1 := firsttask.SubflowId
		status := firsttask.Status
		if strings.EqualFold(string(status), "waiting") {
			nextStepId, err := s.GetStepIdOfEnclosingCallSubflow(flowId, taskname1, strconv.Itoa(subflowid1))
			if err != nil {
				return nil, err
			}
			if nextStepId != "" {
				taskArray, err := s.GetStepdataForActivity(flowId, nextStepId, taskname1)
				if err != nil {
					return nil, err
				}
				stepidInt, _ := strconv.Atoi(stepid)
				taskArray[0].StepId = stepidInt // honour the stepid of starting of subflow
				return taskArray, err
			} else {
				return taskValue, err
			}
		} else {
			return taskValue, err
		}
	} else {
		return taskValue, err
	}
}

func (s *StepStore) GetStepIdOfEnclosingCallSubflow(flowid, taskname, subflowid string) (string, error) {
	set, err := s.db.query("select stepid from steps where taskname = '"+taskname+"' and flowinstanceid= '"+flowid+"' and subflowid= '"+subflowid+"' and status != 'Waiting'", nil)
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetStepIdOfEnclosingCallSubflow after successful connection retry  ")
				set, err = s.db.query("select stepid from steps where taskname = '"+taskname+"' and flowinstanceid= '"+flowid+"' and subflowid= '"+subflowid+"' and status != 'Waiting'", nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return "", err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return "", retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return "", err
		}
	}
	var nextstepid string
	for _, v := range set.Record {
		m := *v
		nextstepid, _ = coerce.ToString(m["stepid"])
	}
	return nextstepid, err
}

func (s *StepStore) GetStepsStatus(flowId string) ([]map[string]string, error) {

	set, err := s.db.query("select stepid, taskname, status, starttime, flowname, rerun, subflowid from steps where flowinstanceid = '"+flowId+"' and stepid != '0' order by cast(stepid as integer)", nil)

	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			if retryErr := s.RetryDBConnection(); retryErr == nil {
				logCache.Debugf("Retrying from GetStepsStatus after successful connection retry  ")
				set, err = s.db.query("select stepid, taskname, status, starttime, flowname, rerun, subflowid from steps where flowinstanceid = '"+flowId+"' and stepid != '0' order by cast(stepid as integer)", nil)
				if err != nil {
					logCache.Errorf("Could not connect to database server error:, %s", err.Error())
					return nil, err
				}
			} else {
				logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
				return nil, retryErr
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", err.Error())
			return nil, err
		}
	}
	var waitingSteps []map[string]string
	var steps []map[string]string
OUTER:
	for _, v := range set.Record {
		m := *v
		stepData := make(map[string]string)
		s1, _ := coerce.ToString(m["stepid"])
		stepData["stepId"] = s1

		status, _ := coerce.ToString(m["status"])
		stepData["status"] = status

		name, _ := coerce.ToString(m["taskname"])
		stepData["taskName"] = name

		flowname, _ := coerce.ToString(m["flowname"])
		stepData["flowname"] = flowname

		rerun, _ := coerce.ToString(m["rerun"])
		stepData["rerun"] = rerun

		subflowid, _ := coerce.ToString(m["subflowid"])
		stepData["subflowid"] = subflowid

		strttime, _ := coerce.ToString(m["starttime"])
		stepData["starttime"] = strttime
		// before appending check if need to merge with step for startasubflow where status=waiting
		// find the pair entry where taskname and subflowid is same and previous status is Waiting
		if strings.EqualFold(status, "completed") || strings.EqualFold(status, "failed") && len(waitingSteps) > 0 {
			for i, waitingStep := range waitingSteps {
				if waitingStep["taskName"] == name && waitingStep["subflowid"] == subflowid {
					// so merge the existing entry as its for corresponding startasubflow
					for _, existingStepdata := range steps {
						if reflect.DeepEqual(existingStepdata, waitingStep) {
							// override value now
							existingStepdata["status"] = status
							waitingSteps = append(waitingSteps[:i], waitingSteps[i+1:]...) // now remove the entry of waiting step as status already merged
							continue OUTER
						}
					}
				}
			}
		}

		steps = append(steps, stepData)
		// add entry into waitingstep for faster check to compare for obly waiting steps
		if strings.EqualFold(status, "waiting") {
			waitingSteps = append(waitingSteps, stepData)
		}
	}
	return steps, err
}

func (s *StepStore) Delete(flowId string) {
	s.Lock()
	delete(s.stepContainers, flowId)
	s.Unlock()
}

type stepContainer struct {
	sync.RWMutex
	status  int
	flowURI string
	steps   []*state.Step
}

func (sc *stepContainer) Status() int {
	sc.RLock()
	status := sc.status
	sc.RUnlock()

	return status
}

func (sc *stepContainer) AddStep(step *state.Step) {
	sc.Lock()

	if len(step.FlowChanges) > 0 {
		if step.FlowChanges[0] != nil && step.FlowChanges[0].SubflowId == 0 {
			if status := step.FlowChanges[0].Status; status != -1 {
				sc.status = status
			}
			if uri := step.FlowChanges[0].FlowURI; uri != "" {
				sc.flowURI = uri
			}
		}
	}

	sc.steps = append(sc.steps, step)
	sc.Unlock()
}

func (sc *stepContainer) Steps() []*state.Step {
	sc.RLock()
	steps := sc.steps
	sc.RUnlock()
	return steps
}

func (s *StepStore) SaveSnapshot(snapshot *state.Snapshot) error {
	//replaces existing snapshot
	return nil
}

func (s *StepStore) GetSnapshot(flowId string) *state.Snapshot {
	if snapshot, ok := s.snapshots.Load(flowId); ok {
		return snapshot.(*state.Snapshot)
	}
	return nil
}

func (s *StepStore) RecordStart(flowState *state.FlowState) error {
	_, err := s.db.InsertFlowState(flowState)
	if err != nil && (err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
		strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
		strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
		strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout")) {
		if retryErr := s.RetryDBConnection(); retryErr == nil {
			logCache.Debug("Retrying from RecordStart after successful connection retry  ")
			_, err = s.db.InsertFlowState(flowState)
			if err != nil {
				logCache.Errorf("Could not connect to database server error:, %s", err.Error())
				return err
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
			return retryErr
		}
	}
	return err
}

func (s *StepStore) RecordEnd(flowState *state.FlowState) error {

	_, err := s.db.UpdateFlowState(flowState)
	if err != nil && (err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
		strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
		strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
		strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout")) {
		if retryErr := s.RetryDBConnection(); retryErr == nil {
			logCache.Debug("Retrying to call  RecordEnd after successful connection retry  ")
			_, err = s.db.UpdateFlowState(flowState)
			if err != nil {
				logCache.Errorf("Could not connect to database server error:, %s", err.Error())
				return err
			}
		} else {
			logCache.Errorf("Could not connect to database server error:, %s", retryErr.Error())
			return retryErr
		}
	}
	return err
}

func (s *StepStore) RetryDBConnection() error {
	conSetting := &pgConnection{}
	err := metadata2.MapToStruct(s.settings, conSetting, false)

	if err != nil {
		logCache.Info("Returning without retry due to connection data parsing issue...")
		return err
	}
	logCache.Info("Trying to ping the database server...")
	dbConnected := 0
	err = s.db.db.Ping()
	// retry attempt on ping only for conn refused and driver bad conn
	if err != nil {
		if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
			strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
			logCache.Info("Failed to ping the database server, trying again...")
			for i := 1; i <= conSetting.MaxConnRetryAttempts; i++ {
				logCache.Infof("Connecting to database server... Attempt-[%d]", i)
				// retry delay
				time.Sleep(time.Duration(conSetting.ConnRetryDelay) * time.Second)
				err = s.db.db.Ping()
				if err != nil {
					if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
						strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
						strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
						strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
						continue
					} else {
						logCache.Errorf("Could not connect to database server %s, %s", conSetting.DbName, err.Error())
						return fmt.Errorf("Could not open connection to database %s, %s", conSetting.DbName, err.Error())
					}
				} else {
					// ping succesful
					dbConnected = 1
					logCache.Infof("Successfully connected to database server in attempt-[%d]", i)
					break
				}
			}
			if dbConnected == 0 {
				logCache.Errorf("Could not connect to database server after %d number of max retry attempts", conSetting.MaxConnRetryAttempts)
				return fmt.Errorf("Could not open connection to database %s, %s", conSetting.DbName, err.Error())
			}
		} else {
			logCache.Errorf("Could not connect to database server %s, %s", conSetting.DbName, err.Error())
			return fmt.Errorf("Could not open connection to database %s, %s", conSetting.DbName, err.Error())
		}
	} else {
		logCache.Info("ping to database server is successful...")
	}
	if dbConnected != 0 {
		logCache.Info("Successfully connected to database server")
	}
	return nil
}
