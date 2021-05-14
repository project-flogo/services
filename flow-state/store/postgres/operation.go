package postgres

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/project-flogo/flow/state"
	"time"
)

const (
	STEP_INSERT      = "INSERT INTO steps (flowinstanceid, stepid, starttime, endtime, stepdata) VALUES ($1,$2,$3,$4,$5);"
	SNAPSHOT_INSERT  = "INSERT INTO snapshopt (flowinstanceid, hostid, stepid, starttime, endtime, stepdata) VALUES ($1,$2,$3,$4,$5,$6);"
	FlowState_INSERT = "INSERT INTO flowstate (flowInstanceId, userId, appId, flowName, hostId,startTime,endTime,status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);"
	UpdateFlowState  = "UPDATE flowstate set endtime=$1,status=$2 where flowinstanceid = $3;"
)

type StatefulDB struct {
	db *sql.DB
}

func (s *StatefulDB) InsertFlowState(flowState *state.FlowState) (results *ResultSet, err error) {
	inputArgs := []interface{}{flowState.FlowInstanceId, flowState.UserId, flowState.AppId, flowState.FlowName, flowState.HostId, flowState.StartTime, flowState.EndTime, flowState.FlowStats}
	return s.insert(FlowState_INSERT, inputArgs)
}

func (s *StatefulDB) UpdateFlowState(flowState *state.FlowState) (results *ResultSet, err error) {
	inputArgs := []interface{}{flowState.EndTime, flowState.FlowStats, flowState.FlowInstanceId}
	return s.insert(UpdateFlowState, inputArgs)
}

func (s *StatefulDB) InsertSteps(step *state.Step) (results *ResultSet, err error) {
	//flowInstanceId := step.FlowInstId
	//hostID := step.HostId
	stepId := step.Id
	//startTime := step.
	b, err := json.Marshal(step)
	if err != nil {
		return nil, err
	}
	stepData := decodeBytes(b)

	inputArgs := []interface{}{step.FlowId, stepId, time.Now(), time.Now(), stepData}
	return s.insert(STEP_INSERT, inputArgs)
}

func (s *StatefulDB) insert(insertSql string, inputArgs []interface{}) (results *ResultSet, err error) {

	// log.Debug("prepared insert [%s]", prepared)
	stmt, err := s.getStepStatement(insertSql)
	if err != nil {
		logCache.Errorf("Failed to prepare statement: %s", err)
		return nil, err
	}

	logCache.Debug("----------- DB Stats in Insert activity -----------")
	rows, err := stmt.Query(inputArgs...)
	if err != nil {
		logCache.Errorf("Executing prepared query got error: %s", err)
		// stmt.Close()
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}
	defer rows.Close()
	return UnmarshalRows(rows)
}

func (s *StatefulDB) query(querySql string, inputArgs []interface{}) (results *ResultSet, err error) {
	// log.Debug("prepared insert [%s]", prepared)
	stmt, err := s.getStepStatement(querySql)
	if err != nil {
		logCache.Errorf("Failed to prepare statement: %s", err)
		return nil, err
	}

	logCache.Debug("----------- DB Stats in query activity -----------")
	rows, err := stmt.Query(inputArgs...)
	if err != nil {
		logCache.Errorf("Executing prepared query got error: %s", err)
		// stmt.Close()
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}
	defer rows.Close()
	return UnmarshalRows(rows)
}

//GetStatement
func (s *StatefulDB) getStepStatement(prepared string) (stmt *sql.Stmt, err error) {
	preparedQueryCacheMutex.Lock()
	defer preparedQueryCacheMutex.Unlock()
	stmt, ok := preparedQueryCache[prepared]
	if !ok {
		stmt, err = s.db.Prepare(prepared)
		if err != nil {
			return nil, err
		}
		preparedQueryCache[prepared] = stmt
	}
	return stmt, nil
}

func decodeBytes(blob []byte) []byte {
	decodedBlob, err := base64.StdEncoding.DecodeString(string(blob))
	if err != nil {
		return blob
	}
	return decodedBlob
}
