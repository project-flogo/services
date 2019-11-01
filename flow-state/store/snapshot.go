package store

import (
	"github.com/project-flogo/flow/state"
)

var snapshotStore SnapshotStore

type SnapshotStore interface {
	GetStatus(flowId string) int
	GetFlow(flowId string) *state.FlowInfo
	GetFlows() []*state.FlowInfo
	SaveSnapshot(snapshot *state.Snapshot) error
	GetSnapshot(flowId string) *state.Snapshot
	Delete(flowId string)
}

func GetSnapshotStore() SnapshotStore {
	return snapshotStore
}

func SetSnapshotStore(store SnapshotStore) {
	snapshotStore = store
}
