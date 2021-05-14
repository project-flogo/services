package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/store/metadata"
	"sync"
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

func (s *StepStore) GetFlow(flowId string) *state.FlowInfo {

	s.RLock()
	sc, ok := s.stepContainers[flowId]
	s.RUnlock()

	if ok {
		return &state.FlowInfo{Id: flowId, Status: sc.Status(), FlowURI: sc.flowURI}
	}

	return nil
}

func (s *StepStore) GetFailedFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error) {
	var whereStr = "where"
	if len(metadata.Username) < 0 {
		return nil, fmt.Errorf("unauthorized, please provide user infomation")
	}
	whereStr += " userId='" + metadata.Username + "'"

	if len(metadata.AppId) < 0 {
		return nil, fmt.Errorf("please provide flow id or flow name")
	}
	whereStr += "  and appId='" + metadata.AppId + "'"

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

func (s *StepStore) GetFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error) {
	var whereStr = "where"
	if len(metadata.Username) < 0 {
		return nil, fmt.Errorf("unauthorized, please provide user infomation")
	}
	whereStr += " userId='" + metadata.Username + "'"

	if len(metadata.AppId) < 0 {
		return nil, fmt.Errorf("please provide flow id or flow name")
	}
	whereStr += "  and appId='" + metadata.AppId + "'"

	if len(metadata.HostId) > 0 {
		whereStr += "  and hostId='" + metadata.HostId + "'"
	}

	if len(metadata.FlowName) > 0 {
		whereStr += "  and flowname='" + metadata.FlowName + "'"
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
		}
		flowinfo = append(flowinfo, info)
	}
	return flowinfo, err
}

func (s *StepStore) SaveStep(step *state.Step) error {
	_, err := s.db.InsertSteps(step)
	return err
}

func (s *StepStore) GetSteps(flowId string) []*state.Step {
	s.RLock()
	sc, ok := s.stepContainers[flowId]
	s.RUnlock()
	if ok {
		return sc.Steps()
	}

	return nil
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
