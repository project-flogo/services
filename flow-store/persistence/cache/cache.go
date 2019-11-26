package cache

import (
	"fmt"
	"github.com/project-flogo/services/flow-store/flow"
	"github.com/project-flogo/services/flow-store/persistence/api"
	"strconv"
	"sync/atomic"
	"time"
)

var flowId int64 = 0

type cacheStorage struct {
	cache *memCache
}

func NewCacheStorage() api.Storage {
	return &cacheStorage{cache: NewCache()}
}

func (f *cacheStorage) AllFlows() []*flow.Flow {
	var flows []*flow.Flow
	for _, v := range f.cache.AllFlows() {
		flows = append(flows, v)
	}
	return flows
}

func (f *cacheStorage) SaveFlow(flow map[string]interface{}) string {
	fl := getFlow(flow)
	f.cache.AddFlow(fl.Metdata.Id, fl)
	return fl.Metdata.Id
}

func (f *cacheStorage) GetFlow(id string) (interface{}, error) {
	fl := f.cache.GetFlow(id)
	if fl == nil {
		return nil, nil
		//return nil, fmt.Errorf("flow [%s] no found", id)
	}
	return fl.Flow, nil
}

func (f *cacheStorage) DeleteFlow(flowId string) error {
	fl := f.cache.GetFlow(flowId)
	if fl == nil {
		return fmt.Errorf("flow [%s] not found", flowId)
	}
	f.cache.DeleteFlow(flowId)
	return nil
}

func (f *cacheStorage) GetFlowMetadata(flowId string) (*flow.Metdata, error) {
	fl := f.cache.GetFlow(flowId)
	if fl == nil {
		return nil, fmt.Errorf("flow [%s] not found", flowId)
	}
	return fl.Metdata, nil
}

func getFlow(body map[string]interface{}) *flow.Flow {

	var id string
	var idtmp, ok = body["id"]
	if ok && idtmp != nil {
		id = idtmp.(string)
	} else {
		id = strconv.FormatInt(incrementalFlowId(), 10)
	}

	var description string
	des, ok := body["description"]
	if ok {
		description = des.(string)
	}
	flowMetadata := &flow.Metdata{
		Id:           id,
		Name:         body["name"].(string),
		Description:  description,
		CreationDate: time.Now().String(),
	}

	flow := &flow.Flow{
		Metdata: flowMetadata,
		Flow:    body,
	}
	return flow
}

func incrementalFlowId() int64 {
	return atomic.AddInt64(&flowId, 1)
}
