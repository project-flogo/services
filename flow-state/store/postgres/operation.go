package postgres

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/project-flogo/flow/state"
	task2 "github.com/project-flogo/services/flow-state/store/task"
	"time"
)

const (
	STEP_INSERT      = "INSERT INTO steps (flowinstanceid, stepid, taskname, status, starttime, endtime, stepdata) VALUES ($1,$2,$3,$4,$5,$6,$7);"
	SNAPSHOT_INSERT  = "INSERT INTO snapshopt (flowinstanceid, hostid, stepid, starttime, endtime, stepdata) VALUES ($1,$2,$3,$4,$5,$6);"
	FlowState_INSERT = "INSERT INTO flowstate (flowInstanceId, userId, appName,appVersion, flowName, hostId,startTime,endTime,status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8, $9);"
	UpdateFlowState  = "UPDATE flowstate set endtime=$1,status=$2, executiontime=(EXTRACT(EPOCH FROM ($1 - starttime)))*1000 where flowinstanceid = $3;"

	UpsertSteps = "INSERT INTO steps (flowinstanceid, stepid, taskname, status, starttime, endtime, stepdata) VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (flowinstanceid, stepid) DO UPDATE SET status = EXCLUDED.status, starttime=EXCLUDED.starttime,endtime= EXCLUDED.endtime,stepdata=EXCLUDED.stepdata;\n"
)

type StatefulDB struct {
	db *sql.DB
}

func (s *StatefulDB) InsertFlowState(flowState *state.FlowState) (results *ResultSet, err error) {
	inputArgs := []interface{}{flowState.FlowInstanceId, flowState.UserId, flowState.AppName, flowState.AppVersion, flowState.FlowName, flowState.HostId, flowState.StartTime, flowState.EndTime, flowState.FlowStats}
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
	tasks, err := task2.StepToTask(step)
	var status, taskName string
	if len(tasks) > 1 {
		for _, task := range tasks {
			taskName = task.Id
			status += "," + string(task.Status)
		}
	} else {
		if len(tasks) == 1 {
			status = string(tasks[0].Status)
			taskName = string(tasks[0].Id)
		}
	}

	//startTime := step.
	b, err := json.Marshal(step)
	if err != nil {
		return nil, err
	}
	stepData := decodeBytes(b)

	inputArgs := []interface{}{step.FlowId, stepId, taskName, status, time.Now(), time.Now(), stepData}
	return s.insert(UpsertSteps, inputArgs)
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

func (s *StatefulDB) update(insertSql string, inputArgs []interface{}) (results *ResultSet, err error) {

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
