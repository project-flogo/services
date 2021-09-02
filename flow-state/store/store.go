package store

import (
	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/store/mem"
	"github.com/project-flogo/services/flow-state/store/metadata"
	"github.com/project-flogo/services/flow-state/store/postgres"
	"os"
)

const (
	Memory     = "Memory"
	File       = "File"
	DynamoDB   = "dynamodb"
	CosmosDB   = "cosmosdb"
	RestServer = "REST"
	Postgres   = "Postgres"

	StoreType = "FLOGO_STATEFUL_STORE_TYPE"
)

type Store interface {
	GetStatus(flowId string) int
	GetFlow(flowId string, metadata *metadata.Metadata) (*state.FlowInfo, error)
	GetFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error)
	GetFailedFlows(metadata *metadata.Metadata) ([]*state.FlowInfo, error)
	SaveStep(step *state.Step) error
	GetSteps(flowId string) ([]*state.Step, error)
	GetStepsNoData(flowId string) ([]map[string]string, error)
	Delete(flowId string)
	SaveSnapshot(snapshot *state.Snapshot) error
	GetSnapshot(flowId string) *state.Snapshot
	RecordStart(step *state.FlowState) error
	RecordEnd(step *state.FlowState) error
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
	switch getStoreType() {
	case Postgres:
		var err error
		store, err = postgres.NewStore(settings)
		if err != nil {
			return err
		}
	case Memory:
		store = mem.NewStore()
	}
	return nil
}

func getStoreType() string {
	v, ok := os.LookupEnv(StoreType)
	if !ok {
		return Postgres
	}
	return v
}