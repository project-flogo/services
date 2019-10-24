package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/julienschmidt/httprouter"
	"github.com/project-flogo/services/flow-state/flow"
	"github.com/project-flogo/services/flow-state/persistence"
)

var Port = flag.String("p", "9190", "The port of the server")
var log = newLogger()

var storage = persistence.GetStorage()

func init() {
	flag.Parse() // get the arguments from command line
}

func newLogger() *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Encoding = "console"
	loggerConfig.DisableCaller = true

	eCfg := loggerConfig.EncoderConfig
	eCfg.TimeKey = "timestamp"
	eCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	eCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	loggerConfig.EncoderConfig = eCfg

	log, err := loggerConfig.Build()
	if err != nil {
		panic("New logger failed," + err.Error())
	}
	return log
}

func main() {
	log.Info("Start flow states go server")
	stateRouter := httprouter.New()
	// TODO remove all non /v1/ apis once web uses versioning
	//Ping
	stateRouter.GET("/v1/ping", Ping)

	//get all steps for the instance
	stateRouter.GET("/v1/instances/:flowID/steps", ListFlowSteps)
	stateRouter.GET("/instances/:flowID/steps", ListFlowSteps)
	stateRouter.GET("/v1/instances/:flowID/status", GetFlowStatus)
	stateRouter.GET("/v1/instances/:flowID", ListFlowSteps)
	stateRouter.DELETE("/v1/instances/:flowID", DeleteFLow)
	stateRouter.GET("/v1/instances/:flowID/snapshot/:snapshotID", GetSnapshotStep)
	stateRouter.GET("/v1/instances/:flowID/metadata", FlowMetadata)
	stateRouter.POST("/v1/instances/snapshot", POSTSnapshot)
	stateRouter.POST("/v1/instances/steps", SaveSteps)

	log.Info("Started server on localhost:" + *Port)
	http.ListenAndServe(":"+*Port, &StateServer{stateRouter})
}

type StateServer struct {
	r *httprouter.Router
}

func (s *StateServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	//Catching unexpected error.
	defer func() {
		if x := recover(); x != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			err, _ := json.Marshal(x)
			log.Error(string(err))
			rw.Write(err)
		}
	}()

	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Atmosphere-Remote-User,X-Atmosphere-Token")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}

	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func Ping(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	fmt.Fprintf(response, "%s", "{\"status\":\"ok\"}")
}

func ListFlowSteps(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowID := params.ByName("flowID")
	log.Info("List flow " + flowID + " status")
	results, err := ListSteps(flowID, false)
	if err != nil {
		ErrorResponse(response, http.StatusInternalServerError, fmt.Errorf("List flow "+flowID+" steps error"))
		log.Error(fmt.Sprintf("List flow "+flowID+" steps error: %v", err))
		return
	} else {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(response).Encode(results); err != nil {
			log.Error(err.Error())
		}
	}
}

func GetFlowStatus(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowID := params.ByName("flowID")
	log.Debug("Get flow " + flowID + " status")
	steps := storage.GetSteps(flowID)
	if len(steps) <= 0 {
		ErrorResponse(response, http.StatusBadRequest, fmt.Errorf("Not flow "+flowID+" found"))
		return
	}
	status := steps[len(steps)-1].Status
	statusM := make(map[string]string)
	statusM["id"] = flowID
	statusM["status"] = strconv.FormatInt(status, 10)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(statusM); err != nil {
		log.Error(err.Error())
	}
}

func DeleteFLow(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowID := params.ByName("flowID")
	log.Debug("Delete flow " + flowID)

	err := storage.DeleteFlow(flowID)
	if err != nil {
		ErrorResponse(response, http.StatusInternalServerError, fmt.Errorf("Get snapshot names error"))
		log.Error(fmt.Sprintf("Get snapshot names error: %v", err))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	//fmt.Fprintf(response, "%d", remSnapShotResp)
}

func ListSteps(flowID string, withStatus bool) (map[string]interface{}, error) {

	var steps []interface{}
	var tasks []interface{}
	taskMetadata := map[string]interface{}{}
	flowMetadata := map[string]interface{}{}
	result := map[string]interface{}{}
	step := map[string]interface{}{}
	taskStates := make(map[string]int64)

	vSteps := storage.GetSteps(flowID)

	for _, change := range vSteps {

		log.Debug(fmt.Sprintf("Change json: %+v", change))

		stepId := change.ID
		stepData := change.StepData
		if stepData == nil {
			log.Debug("No step data found")
		}

		tdchanges := stepData.TdChanges

		snapshot := storage.GetSnapshot(flowID, strconv.FormatInt(stepId, 10))
		snapshotData := snapshot.SnapshotData
		for _, tdchange := range tdchanges {
			taskId := tdchange.ID

			rootTaskEnv := snapshotData.RootTaskEnv
			taskDatas := rootTaskEnv.TaskDatas
			if taskDatas != nil {
				for _, taskData := range taskDatas {
					log.Debug(fmt.Sprintf("taskData :", taskData.TaskId, " taskId: ", taskId))
					if taskId == taskData.TaskId {
						attrs := taskData.Attrs
						log.Debug("Found task match!!!")
						log.Debug(fmt.Sprintf("taskID: ", taskData.TaskId))
						log.Debug(fmt.Sprintf("stepId: ", stepId))
						taskMetadata["taskId"] = taskData.TaskId
						taskMetadata["attributes"] = attrs
						//taskStates[toString(taskData.TaskId)] = taskData.State
						if attrs != nil {
							attributes := attrs.([]interface{})
							if len(attributes) > 0 {
								tasks = append(tasks, taskMetadata)
							}
						}

						taskMetadata = make(map[string]interface{})
					}
				}
			} else {
				log.Debug(fmt.Sprintf("Task datas not found for snapshot: "+flowID+":", stepId))
			}

			taskStates[toString(tdchange.TaskData.TaskId)] = tdchange.TaskData.State
		}

		if snapshotData.ID != "" {
			var tmpTaskId interface{}
			var stepTaskId interface{}
			wqChanges := change.StepData.WqChanges
			if wqChanges != nil && len(wqChanges) > 0 {
				for _, wqchange := range wqChanges {
					tmpTaskId = wqchange.Wtem.TaskId
					chgType := wqchange.ChgType

					if chgType == 3 && !IsRootTask(tmpTaskId) {
						stepTaskId = tmpTaskId
						break
					}
				}

				if stepTaskId == nil {
					stepTaskId = GetRootId(tmpTaskId)
				}
			}

			flowAttrs := snapshotData.Attrs
			if flowAttrs != nil {
				flowMetadata["state"] = snapshotData.State
				flowMetadata["status"] = snapshotData.Status
				flowMetadata["attributes"] = flowAttrs
				step["flow"] = flowMetadata
			}
			step["taskId"] = stepTaskId
			step["id"] = stepId

			if taskState, ok := taskStates[toString(stepTaskId)]; ok {
				step["taskState"] = taskState
			}

			flowMetadata = make(map[string]interface{})
			step["tasks"] = tasks
			steps = append(steps, step)
			tasks = tasks[0:0]
			step = make(map[string]interface{})
			//snapshotDataObj = (flow.SnapshotData{})
		}
	}

	if withStatus {
		status := vSteps[len(steps)-1].Status
		result["status"] = status
	}

	result["steps"] = steps

	return result, nil
}

func toString(val interface{}) string {
	switch t := val.(type) {
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.Itoa(int(t))
	case nil:
		return "___ERROR___"
	default:
		b, err := json.Marshal(t)
		if err != nil {
			fmt.Printf("Unable to Coerce %#v to string\n", t)
			return "___ERROR___"
		}
		return string(b)
	}
}

func GetSnapshotStep(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowId := params.ByName("flowID")
	stepId := params.ByName("snapshotID")

	log.Info("get snapshot step,  flow:" + flowId + " Step id:" + stepId)
	snapshot := storage.GetSnapshot(flowId, stepId)
	if snapshot == nil {
		ErrorResponse(response, http.StatusInternalServerError, fmt.Errorf("No flow "+flowId+" and step "+stepId+" snapshot data found"))
		log.Error(fmt.Sprintf("No flow " + flowId + " and step " + stepId + " snapshot data found"))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(snapshot.SnapshotData); err != nil {
		log.Error(err.Error())
	}
}

func FlowMetadata(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	flowID := params.ByName("flowID")
	log.Info("Get snapshot metadata, flow id: " + flowID)

	metadata := storage.FlowMetadata(flowID)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(metadata); err != nil {
		log.Error(err.Error())
	}
}

func SaveSteps(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		ErrorResponse(response, http.StatusInternalServerError, fmt.Errorf("Read body error"))
		log.Error(fmt.Sprintf("Read body error: %v", err))
		return
	}
	stepInfo := &flow.StepInfo{}
	jsonerr := json.Unmarshal(content, stepInfo)
	if jsonerr != nil {
		ErrorResponse(response, http.StatusInternalServerError, fmt.Errorf("Unmarshal step post body error"))
		log.Debug(fmt.Sprintf("Step content: ", string(content)))
		log.Error(fmt.Sprintf("Unmarshal step post body error %v", jsonerr))
		return
	}

	storage.SaveStep(stepInfo)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}

func POSTSnapshot(response http.ResponseWriter, request *http.Request, params httprouter.Params) {

	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		ErrorResponse(response, http.StatusBadRequest, fmt.Errorf("Read body error"))
		log.Error(fmt.Sprintf("Read body error: %v", err))
		return
	}
	log.Debug(fmt.Sprintf("Snapshot content: ", string(content)))
	snapshot := &flow.Snapshot{}
	jsonerr := json.Unmarshal(content, snapshot)
	if jsonerr != nil {
		ErrorResponse(response, http.StatusInternalServerError, fmt.Errorf("Unmarshal snapshot post body error"))
		log.Debug(fmt.Sprintf("Snapshot content: ", string(content)))
		log.Debug(fmt.Sprintf("Unmarshal snapshot post body error %v", jsonerr))
		return
	}
	storage.SaveSnapshot(snapshot)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}

type StateError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, msg error) *StateError {
	return &StateError{
		Code:    code,
		Message: msg.Error(),
	}
}

func ErrorResponse(response http.ResponseWriter, code int, err error) {
	flowErorr := NewError(code, err)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(code)
	if err := json.NewEncoder(response).Encode(flowErorr); err != nil {
		log.Error(err.Error())
	}
}

func IsRootTask(taskId interface{}) bool {
	strId := ToString(taskId)
	if strId == "1" || strId == "root" {
		return true
	}
	return false
}

func GetRootId(taskId interface{}) interface{} {
	if taskId == nil {
		panic("Invalid nil taskId found")
	}
	switch taskId.(type) {
	case string:
		return "root"
	case float64:
		return 1
	default:
		panic(fmt.Sprintf("Error parsing Task with Id '%v', invalid type '%T'", taskId, taskId))
	}
}

func ToString(m interface{}) string {
	if m == nil {
		panic("Invalid nil activity id found")
	}
	switch m.(type) {
	case string:
		return m.(string)
	case float64, int, int64:
		return strconv.Itoa(int(m.(float64)))
	default:
		v, _ := json.Marshal(m)
		return string(v)
	}
}
