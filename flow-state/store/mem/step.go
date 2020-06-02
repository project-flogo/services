package mem

import (
	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/event"
	"github.com/project-flogo/services/flow-state/store"
	"sync"
)

func init() {
	store.SetStepStore(&StepStore{stepContainers: make(map[string]*stepContainer)})
}

type StepStore struct {
	sync.RWMutex
	stepContainers map[string]*stepContainer
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

func (s *StepStore) GetFlows() []*state.FlowInfo {

	var infos []*state.FlowInfo

	s.RLock()
	for id, value := range s.stepContainers {
		infos = append(infos, &state.FlowInfo{Id: id, FlowURI: value.flowURI, Status: value.Status()})
	}
	s.RUnlock()

	return infos
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
