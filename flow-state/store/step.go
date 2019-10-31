package store

import (
	"github.com/project-flogo/flow/state"
)

var stepStore StepStore

type StepStore interface {
	GetStatus(flowId string) int
	GetFlow(flowId string) *FlowInfo
	GetFlows() []*FlowInfo
	SaveStep(step *state.Step) error
	GetSteps(flowId string) []*state.Step
	Delete(flowId string)
}

//todo add method to get first n steps
//todo potentially have option to construct/update the snapshot as steps come in

func GetStepStore() StepStore {
	return stepStore
}

func SetStepStore(store StepStore) {
	stepStore = store
}

