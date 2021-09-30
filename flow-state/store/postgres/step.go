package postgres

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/store/metadata"
	"github.com/project-flogo/services/flow-state/store/task"
)

func NewStore(settings map[string]interface{}) (*StepStore, error) {
	db, err := NewDB(settings)
	if err != nil {
		return nil, err
	}
	return &StepStore{db: &StatefulDB{db: db}}, err
}

type StepStore struct {
	sync.RWMutex
	db             *StatefulDB
	userId         string
	appId          string
	stepContainers map[string]*stepContainer
	snapshots      sync.Map
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
		return nil, err
	}

	v, _ := json.Marshal(set.Record)
	fmt.Println(string(v))

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
		return nil, err
	}

	v, _ := json.Marshal(set.Record)
	fmt.Println(string(v))

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
	fmt.Println("select flowinstanceid, flowname, status, hostid, starttime, endtime from flowstate " + whereStr)
	set, err := s.db.query("select flowinstanceid, flowname, status, hostid, starttime, endtime from flowstate "+whereStr, nil)
	if err != nil {
		return nil, err
	}

	v, _ := json.Marshal(set.Record)
	fmt.Println(string(v))

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
		whereStr += "  and flowinstanceid='" + mtdata.FlowInstanceId + "'"
	}

	if len(mtdata.Interval) > 0 {
		whereStr += "  and starttime >= NOW() - INTERVAL '" + mtdata.Interval + "'"
	}

	if len(mtdata.Offset) > 0 && len(mtdata.Limit) > 0 {
		offsetLimitStr := "  offset '" + mtdata.Offset + "'  limit  '" + mtdata.Limit + "'"
		whereStr += offsetLimitStr
	}
	fmt.Println("select flowinstanceid, flowname, status, hostid, starttime, endtime, executiontime, count(*) over() AS full_count from flowstate " + whereStr)
	set, err := s.db.query("select flowinstanceid, flowname, status, hostid, starttime, endtime, executiontime, count(*) over() AS full_count from flowstate "+whereStr, nil)
	if err != nil {
		return nil, err
	}

	v, _ := json.Marshal(set.Record)
	fmt.Println(string(v))
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
		info := &state.FlowInfo{
			Id:            id,
			FlowName:      flowName,
			HostId:        hostid,
			FlowStatus:    status,
			StartTime:     starttime,
			EndTime:       endtime,
			ExecutionTime: executiontime,
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

	set, err := s.db.query("select flowinstanceid, flowname, status from flowstate "+whereStr, nil)
	if err != nil {
		return nil, err
	}

	v, _ := json.Marshal(set.Record)
	fmt.Println(string(v))

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
			FlowURI:    "res://flow:" + flowName,
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
		return nil, err
	}

	var flownameArray []string
	for _, v := range set.Record {
		m := *v
		flowName, _ := coerce.ToString(m["flowname"])
		flownameArray = append(flownameArray, flowName)
	}
	return flownameArray, err
}

func (s *StepStore) SaveStep(step *state.Step) error {
	_, err := s.db.InsertSteps(step)
	return err
}

func (s *StepStore) GetSteps(flowId string) ([]*state.Step, error) {

	set, err := s.db.query("select stepdata from steps where flowinstanceid = '"+flowId+"'", nil)
	if err != nil {
		return nil, err
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
	return steps, err
}

func (s *StepStore) GetStepsAsTasks(flowId string) ([][]*task.Task, error) {

	set, err := s.db.query("select stepdata from steps where flowinstanceid = '"+flowId+"'", nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	var step *state.Step
	for _, v := range set.Record {
		m := *v
		s1, err := coerce.ToBytes(m["stepdata"])
		if err != nil {
			fmt.Println("error", err)
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
	return taskValue, err
}

func (s *StepStore) GetStepsStatus(flowId string) ([]map[string]string, error) {

	set, err := s.db.query("select stepid, taskname, status, starttime from steps where flowinstanceid = '"+flowId+"'", nil)
	if err != nil {
		return nil, err
	}

	var steps []map[string]string
	for _, v := range set.Record {
		m := *v
		stepData := make(map[string]string)
		s1, _ := coerce.ToString(m["stepid"])
		stepData["stepId"] = s1

		status, _ := coerce.ToString(m["status"])
		stepData["status"] = status

		name, _ := coerce.ToString(m["taskname"])
		stepData["taskName"] = name

		strttime, _ := coerce.ToString(m["starttime"])
		stepData["starttime"] = strttime
		steps = append(steps, stepData)
	}

	if len(steps) <= 0 {
		return nil, fmt.Errorf("step for flow instance [%s] not found", flowId)
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
	return err
}

func (s *StepStore) RecordEnd(flowState *state.FlowState) error {
	_, err := s.db.UpdateFlowState(flowState)
	return err
}
