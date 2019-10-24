package api

import "github.com/project-flogo/services/flow-state/flow"

type Storage interface {
	SaveStep(step *flow.StepInfo)
	SaveSnapshot(snapshot *flow.Snapshot)
	GetSnapshot(flowId, stepID string) *flow.Snapshot
	ListFlowSteps(flwoid string) *flow.StepInfo
	GetSteps(flowID string) []*flow.StepInfo
	GetSnapshotStatus(flwoid string) map[string]string
	GetSnapshotSteps(flwoid string) map[string]interface{}
	FlowMetadata(flwoid string) map[string]interface{}
	DeleteFlow(path string) error
}
