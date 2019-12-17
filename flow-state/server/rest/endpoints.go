package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/store"
)

type ServiceEndpoints struct {
	logger        log.Logger
	stepStore     store.StepStore
	snapshotStore store.SnapshotStore
}

func AppendEndpoints(router *httprouter.Router, logger log.Logger, exposeRecorder bool) {

	sm := &ServiceEndpoints{
		logger:        logger,
		stepStore:     store.GetStepStore(),
		snapshotStore: store.GetSnapshotStore(),
	}

	router.GET("/v1/instances", sm.getInstances)
	router.GET("/v1/instances/:flowId/details", sm.getInstance)
	router.GET("/v1/instances/:flowId/status", sm.getStatus)

	router.GET("/v1/instances/:flowId/steps", sm.getSteps)
	router.GET("/v1/instances/:flowId/snapshot", sm.getSnapshot)
	router.GET("/v1/instances/:flowId/snapshot/:stepId", sm.getSnapshotAtStep)
	router.DELETE("/v1/instances/:flowId", sm.deleteInstance)

	if exposeRecorder {
		router.POST("/v1/instances/snapshot", sm.saveSnapshot)
		router.POST("/v1/instances/steps", sm.saveStep)
	}
}

func (se *ServiceEndpoints) getInstances(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	se.logger.Debugf("Endpoint[GET:/instances] : Called")
	instances := se.snapshotStore.GetFlows()

	if len(instances) == 0 {
		se.logger.Debugf("Getting instances from steps")
		instances = se.stepStore.GetFlows()
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	if instances == nil {
		_, _ = response.Write([]byte("[]"))
		return
	}

	if err := json.NewEncoder(response).Encode(instances); err != nil {
		se.logger.Error(err.Error())
	}
}

func (se *ServiceEndpoints) getInstance(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowId := params.ByName("flowId")
	se.logger.Debugf("Endpoint[GET:/instances/%s] : Called", flowId)
	instance := se.snapshotStore.GetFlow(flowId)

	if instance == nil {
		se.logger.Debugf("Getting instance from steps")
		instance = se.stepStore.GetFlow(flowId)

		if instance == nil {
			response.WriteHeader(http.StatusNotFound)
			return
		}
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(instance); err != nil {
		se.logger.Error(err.Error())
	}
}

func (se *ServiceEndpoints) getStatus(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowId := params.ByName("flowId")
	se.logger.Debugf("Endpoint[GET:/instances/%s/status] : Called", flowId)
	status := se.snapshotStore.GetStatus(flowId)

	if status == -1 {
		se.logger.Debugf("Getting status from steps")
		status = se.stepStore.GetStatus(flowId)

		if status == -1 {
			response.WriteHeader(http.StatusNotFound)
			return
		}
	}

	statusObj := make(map[string]string, 1)
	statusObj["status"] = strconv.Itoa(status)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(statusObj); err != nil {
		se.logger.Error(err.Error())
	}
}

func (se *ServiceEndpoints) getSteps(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowId := params.ByName("flowId")
	se.logger.Debugf("Endpoint[GET:/instances/%s/steps] : Called", flowId)
	steps := se.stepStore.GetSteps(flowId)

	if steps == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(steps); err != nil {
		se.logger.Error(err.Error())
	}
}

func (se *ServiceEndpoints) getSnapshot(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowId := params.ByName("flowId")
	se.logger.Debugf("Endpoint[GET:/instances/%s/snapshot] : Called", flowId)
	snapshot := se.snapshotStore.GetSnapshot(flowId)

	if snapshot == nil {
		se.logger.Debugf("Getting Snapshot from steps")
		steps := se.stepStore.GetSteps(flowId)

		if steps == nil {
			response.WriteHeader(http.StatusNotFound)
			return
		}

		snapshot = state.StepsToSnapshot(flowId, steps)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(snapshot); err != nil {
		se.logger.Error(err.Error())
	}
}

func (se *ServiceEndpoints) getSnapshotAtStep(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowId := params.ByName("flowId")
	stepIdStr := params.ByName("stepId")

	se.logger.Debugf("Endpoint[GET:/instances/%s/snapshot/%s] : Called", flowId, stepIdStr)
	steps := se.stepStore.GetSteps(flowId)

	if steps == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	stepId, err := strconv.Atoi(stepIdStr)
	if err != nil {
		se.error(response, http.StatusBadRequest, fmt.Errorf("invalid stepId: %s", stepIdStr))
		se.logger.Errorf("Endpoint[GET:/instances/%s/snapshot/%s] : Invalid StepId")
		return
	}

	if stepId >= len(steps) {
		se.error(response, http.StatusBadRequest, fmt.Errorf("invalid stepId: %d, only %d exists", stepId, len(steps)))
		se.logger.Errorf("Endpoint[GET:/instances/%s/snapshot/%s] : Step does not exists")
		return
	}

	snapshot := state.StepsToSnapshot(flowId, steps[:stepId+1])

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(snapshot); err != nil {
		se.logger.Error(err.Error())
	}
}

func (se *ServiceEndpoints) deleteInstance(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowId := params.ByName("flowId")
	se.logger.Debugf("Endpoint[DEL:/instances/%s] : Called", flowId)

	se.snapshotStore.Delete(flowId)
	se.stepStore.Delete(flowId)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}

func (se *ServiceEndpoints) saveStep(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	se.logger.Debugf("Endpoint[POST:/instances/steps] : Called")

	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		se.error(response, http.StatusBadRequest, fmt.Errorf("unable to read body"))
		se.logger.Error("Endpoint[POST:/instances/steps] : %v", err)
		return
	}

	step := &state.Step{}
	err = json.Unmarshal(content, step)
	if err != nil {
		se.error(response, http.StatusBadRequest, fmt.Errorf("unable to unmarshal step json"))
		se.logger.Debugf("Endpoint[POST:/instances/steps] : Step content - %s ", string(content))
		se.logger.Errorf("Endpoint[POST:/instances/steps] : Error unmarshalling step - %v", err)
		return
	}

	err = se.stepStore.SaveStep(step)
	if err != nil {
		se.error(response, http.StatusInternalServerError, fmt.Errorf("unable to save step"))
		se.logger.Errorf("Endpoint[POST:/instances/steps] : Error saving step - %v", err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}

func (se *ServiceEndpoints) saveSnapshot(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	se.logger.Debugf("Endpoint[POST:/instances/snapshot] : Called")

	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		se.error(response, http.StatusBadRequest, fmt.Errorf("unable to read body"))
		se.logger.Error("Endpoint[POST:/instances/snapshot] : %v", err)
		return
	}

	snapshot := &state.Snapshot{SnapshotBase: &state.SnapshotBase{}}
	err = json.Unmarshal(content, snapshot)
	if err != nil {
		se.error(response, http.StatusBadRequest, fmt.Errorf("unable to unmarshal snapshot json"))
		se.logger.Debugf("Endpoint[POST:/instances/snapshot] : Snapshot content - %s ", string(content))
		se.logger.Errorf("Endpoint[POST:/instances/snapshot] : Error unmarshalling snapshot - %v", err)
		return
	}

	err = se.snapshotStore.SaveSnapshot(snapshot)
	if err != nil {
		se.error(response, http.StatusInternalServerError, fmt.Errorf("unable to save snapshot"))
		se.logger.Errorf("Endpoint[POST:/instances/snapshot] : Error saving snapshot - %v", err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}

type StateError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (se *ServiceEndpoints) error(response http.ResponseWriter, code int, err error) {
	flowError := &StateError{
		Code:    code,
		Message: err.Error(),
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(code)
	if err := json.NewEncoder(response).Encode(flowError); err != nil {
		se.logger.Errorf("unable to encode err to json: %v", err)
	}
}
