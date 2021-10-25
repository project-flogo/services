package mem

import (
	"github.com/project-flogo/services/flow-state/store/mem"
	"sync"

	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/store"
)

func init() {
	store.SetSnapshotStore(&DynaSnapshotStore{})
}

type DynaSnapshotStore struct {
}

func (s *DynaSnapshotStore) GetStatus(flowId string) int {
	if snapshot, ok := s.snapshots.Load(flowId); ok {
		fs := snapshot.(*state.Snapshot)
		return fs.Status
	}
	return -1
}

func (s *DynaSnapshotStore) GetFlow(flowId string) *state.FlowInfo {
	if snapshot, ok := s.snapshots.Load(flowId); ok {
		fs := snapshot.(*state.Snapshot)
		return &state.FlowInfo{Id: fs.Id, Status: fs.Status, FlowURI: fs.FlowURI}
	}
	return nil
}

func (s *DynaSnapshotStore) GetFlows() []*state.FlowInfo {

	var infos []*state.FlowInfo

	s.snapshots.Range(func(k, v interface{}) bool {
		fs := v.(*state.Snapshot)
		infos = append(infos, &state.FlowInfo{Id: fs.Id, Status: fs.Status, FlowURI: fs.FlowURI})
		return true
	})

	return infos
}

func (s *DynaSnapshotStore) SaveSnapshot(snapshot *state.Snapshot) error {
	//replaces existing snapshot
	s.snapshots.Store(snapshot.Id, snapshot)
	return nil
}

func (s *DynaSnapshotStore) GetSnapshot(flowId string) *state.Snapshot {
	if snapshot, ok := s.snapshots.Load(flowId); ok {
		return snapshot.(*state.Snapshot)
	}
	return nil
}

func (s *DynaSnapshotStore) Delete(flowId string) {
	s.snapshots.Delete(flowId)
}
