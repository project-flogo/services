package cache

import (
	"strconv"
	"strings"

	"github.com/project-flogo/services/flow-state/flow"
	"github.com/project-flogo/services/flow-state/persistence/api"
)

var flowId int64 = 0

type cacheStorage struct {
	cache *memCache
}

func NewCacheStorage() api.Storage {
	return &cacheStorage{cache: NewCache()}
}

func (f *cacheStorage) SaveStep(steps *flow.StepInfo) {
	f.cache.AddSteps(steps.FlowID, steps)
}

func (f *cacheStorage) SaveSnapshot(snapshot *flow.Snapshot) {
	f.cache.AddSnapshots(snapshot.FlowID+":"+strconv.FormatInt(snapshot.ID, 10), snapshot)
}

func (f *cacheStorage) GetSnapshot(flowId, stepId string) *flow.Snapshot {
	key := flowId + ":" + stepId
	return f.cache.snapshots[key]
}

func (f *cacheStorage) GetSteps(id string) []*flow.StepInfo {
	return f.cache.GetSteps(id)
}

func (f *cacheStorage) GetSnapshotStatus(flowId string) map[string]string {
	return nil
}

func (f *cacheStorage) GetSnapshotSteps(flowId string) map[string]interface{} {
	return nil
}

func (f *cacheStorage) ListFlowSteps(flowId string) *flow.StepInfo {
	return nil
}

func (f *cacheStorage) FlowMetadata(flowId string) map[string]interface{} {
	maxStespId := 0
	snapshots := f.cache.GetAllSnapshots()
	for k, _ := range snapshots {
		keys := strings.Split(k, ":")
		fId, sId := keys[0], keys[1]
		if strings.TrimSpace(fId) == flowId {
			stepId, _ := strconv.Atoi(sId)
			if stepId > maxStespId {
				maxStespId = stepId
			}
		}
	}

	s := f.cache.GetSnapshots(flowId + ":" + strconv.Itoa(maxStespId))
	metadata := make(map[string]interface{})
	metadata["id"] = s.ID
	metadata["flowID"] = s.FlowID
	metadata["status"] = s.Status
	metadata["state"] = s.State
	metadata["date"] = s.Date
	return metadata
}

func (f *cacheStorage) DeleteFlow(flowId string) error {
	return nil
}
