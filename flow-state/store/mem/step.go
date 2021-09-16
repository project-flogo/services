package mem

import (
	"sync"

	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/event"
	"github.com/project-flogo/services/flow-state/store/metadata"
	"github.com/project-flogo/services/flow-state/store/task"
)

//func init() {
//	store.SetStepStore(&StepStore{stepContainers: make(map[string]*stepContainer)})
//}

func NewStore() *StepStore {
	return &StepStore{stepContainers: make(map[string]*stepContainer)}
}

type StepStore struct {
	sync.RWMutex
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

func (s *StepStore) GetFlow(flowid string, fmetadata *metadata.Metadata) (*state.FlowInfo, error) {

	s.RLock()
	sc, ok := s.stepContainers[flowid]
	s.RUnlock()

	if ok {
		return &state.FlowInfo{Id: flowid, Status: sc.Status(), FlowURI: sc.flowURI}, nil
	}

	return nil, nil
}

func (s *StepStore) GetFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error) {

	var infos []*state.FlowInfo

	s.RLock()
	for id, value := range s.stepContainers {
		infos = append(infos, &state.FlowInfo{Id: id, FlowURI: value.flowURI, Status: value.Status()})
	}
	s.RUnlock()

	return infos, nil
}

func (s *StepStore) GetFlowsWithRecordCount(metadata *metadata.Metadata) (*metadata.FlowRecord, error) {

	var infos []*state.FlowInfo

	s.RLock()
	for id, value := range s.stepContainers {
		infos = append(infos, &state.FlowInfo{Id: id, FlowURI: value.flowURI, Status: value.Status()})
	}
	s.RUnlock()

	return nil, nil
}

func (s *StepStore) GetFailedFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error) {

	var infos []*state.FlowInfo

	s.RLock()
	for id, value := range s.stepContainers {
		if value.Status() == 500 {
			infos = append(infos, &state.FlowInfo{Id: id, FlowURI: value.flowURI, Status: value.Status()})
		}
	}
	s.RUnlock()

	return infos, nil
}

func (s *StepStore) GetCompletedFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error) {

	var infos []*state.FlowInfo

	s.RLock()
	for id, value := range s.stepContainers {
		if value.Status() == 100 {
			infos = append(infos, &state.FlowInfo{Id: id, FlowURI: value.flowURI, Status: value.Status()})
		}
	}
	s.RUnlock()

	return infos, nil
}

func (s *StepStore) SaveStep(step *state.Step) error {
	event.PostStepEvent(step)
	s.RLock()
	sc, ok := s.stepContainers[step.FlowId]
	s.RUnlock()

	if !ok {
		s.Lock()
		sc, ok = s.stepContainers[step.FlowId]
		if !ok {
			sc = &stepContainer{}
		}
		s.stepContainers[step.FlowId] = sc
		s.Unlock()
	}

	sc.AddStep(step)

	return nil
}

func (s *StepStore) GetSteps(flowId string) ([]*state.Step, error) {
	s.RLock()
	sc, ok := s.stepContainers[flowId]
	s.RUnlock()
	if ok {
		return sc.Steps(), nil
	}

	return nil, nil
}

func (s *StepStore) GetStepsAsTasks(flowId string) ([][]*task.Task, error) {

	return nil, nil
}
func (s *StepStore) GetStepsNoData(flowId string) ([]map[string]string, error) {
	s.RLock()
	_, ok := s.stepContainers[flowId]
	s.RUnlock()
	if ok {
		return nil, nil
	}

	return nil, nil
}

func (s *StepStore) GetStepdataForActivity(flowId, stepid, taskname string) ([]*task.Task, error) {
	return nil, nil
}
func (s *StepStore) GetFlowNames(metadata *metadata.Metadata) ([]string, error) {
	return nil, nil
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
	s.snapshots.Store(snapshot.Id, snapshot)
	return nil
}

func (s *StepStore) GetSnapshot(flowId string) *state.Snapshot {
	if snapshot, ok := s.snapshots.Load(flowId); ok {
		return snapshot.(*state.Snapshot)
	}
	return nil
}

func (s *StepStore) RecordStart(flowState *state.FlowState) error {
	return nil
}

func (s *StepStore) RecordEnd(flowState *state.FlowState) error {
	return nil
}
