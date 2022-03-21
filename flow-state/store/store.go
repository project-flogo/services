package store

import (
	"fmt"
	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/store/mem"
	"github.com/project-flogo/services/flow-state/store/metadata"
	"github.com/project-flogo/services/flow-state/store/postgres"
	"github.com/project-flogo/services/flow-state/store/task"
)

const (
	Memory     = "memory"
	File       = "File"
	DynamoDB   = "dynamodb"
	CosmosDB   = "cosmosdb"
	RestServer = "REST"
	Postgres   = "postgres"
)

type Store interface {
	GetStatus(flowId string) int
	GetFlow(flowId string, metadata *metadata.Metadata) (*state.FlowInfo, error)
	GetFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error)
	GetFailedFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error)
	GetCompletedFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error)
	GetFlowsWithRecordCount(metadata *metadata.Metadata) (*metadata.FlowRecord, error)
	SaveStep(step *state.Step) error
	GetSteps(flowId string) ([]*state.Step, error)
	GetStepsAsTasks(flowId string) ([][]*task.Task, error)
	GetStepsStatus(flowId string) ([]map[string]string, error)
	GetStepdataForActivity(flowId, stepid, taskname string) ([]*task.Task, error)
	GetFlowNames(metadata *metadata.Metadata) ([]string, error)
	GetAppVersions(metadata *metadata.Metadata) ([]string, error)
	GetAppState(metadata *metadata.Metadata) (string, error)
	SaveAppState(metadata *metadata.Metadata) error
	Delete(flowId string)
	SaveSnapshot(snapshot *state.Snapshot) error
	GetSnapshot(flowId string) *state.Snapshot
	RecordStart(step *state.FlowState) error
	RecordEnd(step *state.FlowState) error
	DeleteSteps(flowId string, stepId string) error
}

//type SnapshotStore interface {
//	GetStatus(flowId string) int
//	GetFlow(flowId string) *state.FlowInfo
//	GetFlows() []*state.FlowInfo
//	SaveSnapshot(snapshot *state.Snapshot) error
//	GetSnapshot(flowId string) *state.Snapshot
//	Delete(flowId string)
//}

var store Store

func RegistedStore() Store {
	return store
}

func InitStorage(settings map[string]interface{}) error {

	if len(settings) == 0 {
		//Default set to mem
		store = mem.NewStore()
		return nil
	}

	persistenceType := settings["type"]
	switch persistenceType {
	case Postgres:
		fmt.Println("Store type is: Postgres")
		var err error
		store, err = postgres.NewStore(settings)
		if err != nil {
			return err
		}
	case Memory:
		fmt.Println("Store type is: Memory")

		store = mem.NewStore()
	}
	return nil
}
