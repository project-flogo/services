package store

import (
	"github.com/project-flogo/flow/state"
)

var snapshotStore SnapshotStore

type SnapshotStore interface {
	GetStatus(flowId string) int
	GetFlow(flowId string) *FlowInfo
	GetFlows() []*FlowInfo
	SaveSnapshot(snapshot *state.Snapshot) error
	GetSnapshot(flowId string) *state.Snapshot
	Delete(flowId string)
}

type FlowInfo struct {
	Id      string `json:"id"`
	FlowURI string `json:"flowURI"`
	Status  int    `json:"status"`
}

func GetSnapshotStore() SnapshotStore {
	return snapshotStore
}

func SetSnapshotStore(store SnapshotStore) {
	snapshotStore = store
}
