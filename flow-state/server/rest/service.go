package rest

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"strings"

	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
)

const (
	CorsPrefix = "STATE_SERVICE"

	SettingPort           = "port"
	SettingExposeRecorder = "exposeRecorder"
	SettingEnableTLS      = "enableTLS"
	SettingCertFile       = "certFile"
	SettingKeyFile        = "keyFile"
)

// StateService is an implementation of StateService service
// that can access flows via URI
type StateService struct {
	server  *Server
	logger  log.Logger
	enabled bool
}

// NewRemoteStateService creates a new StateService
func NewStateService(config *support.ServiceConfig, logger log.Logger) *StateService {

	recorder := &StateService{enabled: config.Enabled, logger:logger}
	_ = recorder.init(config.Settings)
	//todo handle error

	return recorder
}

func (ss *StateService) Name() string {
	return "FlowStateService"
}

func (ss *StateService) Enabled() bool {
	return ss.enabled
}

// Start implements util.Managed.Start()
func (ss *StateService) Start() error {
	return ss.server.Start()
}

// Stop implements util.Managed.Stop()
func (ss *StateService) Stop() error {
	return ss.server.Stop()
}

// Init implements services.StateServiceService.Init()
func (ss *StateService) init(settings map[string]string) error {

	port, set := settings[SettingPort]
	if !set {
		return fmt.Errorf("FlowStateService: required setting 'port' not set")
	}

	router := httprouter.New()

	exposeRecorder := false
	if expose, set := settings[SettingExposeRecorder]; set {
		exposeRecorder = strings.EqualFold("true", expose)
	}

	AppendEndpoints(router, ss.logger, true, exposeRecorder)

	addr := ":" + strings.TrimSpace(port)

	var options []func(*Server)

	enableTLS := false
	if strEnableTLS, set := settings[SettingEnableTLS]; set {
		enableTLS = strings.EqualFold("true", strEnableTLS)
	}

	if enableTLS {
		certFile, _ := settings[SettingCertFile]
		keyFile, _ := settings[SettingKeyFile]
		options = append(options, TLS(certFile, keyFile))
	}

	options = append(options, Logger(ss.logger))

	server, err := newServer(addr, router, options...)
	if err != nil {
		return err
	}

	ss.server = server

	return nil
}
