package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/flow/service"
	"github.com/project-flogo/flow/state"
)

// StateRecorder is an implementation of StateRecorder service
// that can access flows via URI
type StateRecorder struct {
	host    string
	enabled bool
	logger  log.Logger
}

// NewRemoteStateRecorder creates a new StateRecorder
func NewStateRecorder(config *support.ServiceConfig) *StateRecorder {

	recorder := &StateRecorder{enabled: config.Enabled}
	recorder.init(config.Settings)

	//todo switch this logger
	recorder.logger = log.RootLogger()

	return recorder
}

func (sr *StateRecorder) Name() string {
	return service.ServiceStateRecorder
}

func (sr *StateRecorder) Enabled() bool {
	return sr.enabled
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
func (sr *StateRecorder) init(settings map[string]string) {

	host, set := settings["host"]
	port, set := settings["port"]

	if !set {
		panic("StateRecorder: required setting 'host' not set")
	}

	if strings.Index(host, "http") != 0 {
		sr.host = "http://" + host + ":" + port
	} else {
		sr.host = host + ":" + port
	}

	sr.logger.Debugf("StateRecorder: StateRecorder Server = %s", sr.host)
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

func DefaultConfig() *support.ServiceConfig {
	return &support.ServiceConfig{Name: service.ServiceStateRecorder, Enabled: true, Settings: map[string]string{"host": ""}}
}
