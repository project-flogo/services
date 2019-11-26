package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/project-flogo/services/flow-store/flow"
	"github.com/project-flogo/services/flow-store/persistence"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"net/http"
)

var log = newLogger()

var Port = flag.String("p", "9090", "The port of the server")

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
	log.Info("Start flow go server")
	flowRouter := httprouter.New()
	flowRouter.GET("/ping", Ping)

	//New Apis V1
	flowRouter.GET("/v1/flows", ListAllFlow)
	flowRouter.GET("/v1/flows/:id", GetFlow)
	flowRouter.GET("/v1/flows/:id/metadata", GetFlowMetadata)
	flowRouter.POST("/v1/flows", SaveFlow)
	flowRouter.DELETE("/v1/flows/:id", DeleteFlow)

	log.Info("Started server on localhost:" + *Port)
	http.ListenAndServe(":"+*Port, &FlowServer{flowRouter})
}

type FlowServer struct {
	r *httprouter.Router
}

func (s *FlowServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

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
	//log.Info("ping.....")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	fmt.Fprintf(response, "%s", "{\"status\":\"ok\"}")
}

func ListAllFlow(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	log.Debug("List all flows")
	flows := storage.AllFlows()
	var metdatas []*flow.Metdata
	for _, v := range flows {
		metdatas = append(metdatas, v.Metdata)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(metdatas); err != nil {
		log.Error(err.Error())
	}

}

func GetFlow(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	log.Debug("Get flow " + id)

	fl, err := storage.GetFlow(id)
	if err != nil {
		getErr := fmt.Errorf("error getting flow [%s]: %s", id, err)
		handlerErrorResponse(response, http.StatusInternalServerError, getErr)
		log.Error(getErr.Error())
		return
	} else if fl == nil {
		handlerErrorResponse(response, http.StatusNotFound, fmt.Errorf("flow [%s] not found", id))
	} else {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(response).Encode(fl); err != nil {
			log.Error(err.Error())
		}
	}
}

func GetFlowMetadata(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	log.Debug("Get flow metadata" + id)
	fl, err := storage.GetFlowMetadata(id)
	if err != nil {
		handlerErrorResponse(response, http.StatusInternalServerError, fmt.Errorf("Get flow metadata [%s] error [%s]", id, err))
		log.Error(fmt.Sprintf("Get flow metadata "+id+" error :%v", err))
		return
	} else {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(response).Encode(fl); err != nil {
			log.Error(err.Error())
		}
	}
}

func SaveFlow(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	log.Debug("Running save flow.")
	flowContent, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handlerErrorResponse(response, http.StatusBadRequest, err)
		return
	}

	flowInfo := map[string]interface{}{}
	unmarshalErr := json.Unmarshal(flowContent, &flowInfo)
	if unmarshalErr != nil {
		handlerErrorResponse(response, http.StatusInternalServerError, fmt.Errorf("Unmarshal flow body error while save flow"))
		log.Error(fmt.Sprintf("Unmarshal flow body error while save flow:%v", unmarshalErr))
		return
	}

	id := storage.SaveFlow(flowInfo)
	metadata, err := storage.GetFlowMetadata(id)
	if err != nil {
		handlerErrorResponse(response, http.StatusInternalServerError, fmt.Errorf("Get flow from BD error, flow id: "+id))
		log.Error(fmt.Sprintf("Get flow from BD error, flow id: "+id+" :%v", err))
		return
	} else {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(response).Encode(metadata); err != nil {
			log.Error(err.Error())
		}
	}
}

func DeleteFlow(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	log.Info("Delete flow " + id)

	err := storage.DeleteFlow(id)
	if err != nil {
		handlerErrorResponse(response, http.StatusInternalServerError, fmt.Errorf("Delete flow [%s] error [%s]", id, err.Error()))
		log.Error(fmt.Sprintf("Delete flow "+id+" error :%v", err))
		return
	} else {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
	}
}

func handlerErrorResponse(response http.ResponseWriter, code int, err error) {
	flowErorr := NewError(err, code)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(code)
	if err := json.NewEncoder(response).Encode(flowErorr); err != nil {
		log.Error(err.Error())
	}
}

func NewError(err error, code int) *FlowError {
	return &FlowError{
		Code:    code,
		Message: err.Error(),
	}
}

type FlowError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
