package api

import "github.com/project-flogo/services/flow-store/flow"

type Storage interface {
	AllFlows() []*flow.Flow
	SaveFlow(flow map[string]interface{}) string
	GetFlow(path string) (interface{}, error)
	DeleteFlow(path string) error
	GetFlowMetadata(flowId string) (*flow.Metadata, error)
}
