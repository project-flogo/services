package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/service"
	"github.com/project-flogo/flow/state"
)

func init() {
	_ = service.RegisterFactory(&StateRecorderFactory{})
}

type StateRecorderFactory struct {

}

func (s StateRecorderFactory) NewService(config *service.Config) (service.Service, error) {
	recorder := &StateRecorder{}
	err := recorder.init(config.Settings)
	if err != nil {
		return nil, err
	}

	//todo switch this logger
	recorder.logger = log.RootLogger()

	return recorder, nil
}

// StateRecorder is an implementation of StateRecorder service
// that can access flows via URI
type StateRecorder struct {
	host    string
	logger  log.Logger
}

func (sr *StateRecorder) Name() string {
	return "FlowStateRecorder"
}

// Start implements util.Managed.Start()
func (sr *StateRecorder) Start() error {
	// no-op
	return nil
}

// Stop implements util.Managed.Stop()
func (sr *StateRecorder) Stop() error {
	// no-op
	return nil
}

// Init implements services.StateRecorderService.Init()
func (sr *StateRecorder) init(settings map[string]interface{}) error {

	sHost, set := settings["host"]
	if !set {
		return fmt.Errorf("StateRecorder: required setting 'host' not set")
	}
	host, err := coerce.ToString(sHost)
	if err != nil {
		return fmt.Errorf("StateRecorder: invalid host '%v'", sHost)
	}

	sPort, set := settings["port"]
	if !set {
		return fmt.Errorf("StateRecorder: required setting 'port' not set")
	}
	port, err := coerce.ToInt(sPort)
	if err != nil {
		return fmt.Errorf("StateRecorder: invalid port '%v'", sPort)
	}

	if strings.Index(host, "http") != 0 {
		sr.host = "http://" + host + ":" + strconv.Itoa(port)
	} else {
		sr.host = host + ":" + strconv.Itoa(port)
	}

	sr.logger.Debugf("StateRecorder: StateRecorder Server = %s", sr.host)
	return nil
}

// RecordSnapshot implements instance.StateRecorder.RecordSnapshot
func (sr *StateRecorder) RecordSnapshot(snapshot *state.Snapshot) error {

	uri := sr.host + "/v1/instances/snapshot"

	sr.logger.Debugf("POST Snapshot: %s\n", uri)

	jsonReq, _ := json.Marshal(snapshot)

	sr.logger.Debug("JSON: ", string(jsonReq))

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sr.logger.Debug("response Status:", resp.Status)

	if resp.StatusCode >= 300 {
		//todo return error
	}

	return nil
}

// RecordStep implements instance.StateRecorder.RecordStep
func (sr *StateRecorder) RecordStep(step *state.Step) error {

	uri := sr.host + "/v1/instances/steps"

	sr.logger.Debugf("POST Step: %s\n", uri)

	jsonReq, _ := json.Marshal(step)

	sr.logger.Debug("JSON: ", string(jsonReq))

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sr.logger.Debug("response Status:", resp.Status)

	if resp.StatusCode >= 300 {
		//todo return error
	}

	return nil
}
