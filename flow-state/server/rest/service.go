package rest

import (
	"fmt"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/service"
	"github.com/rs/cors"
)

const (
	SettingPort           = "port"
	SettingExposeRecorder = "exposeRecorder"
	SettingStreamingStep  = "streamingStep"
	SettingEnableTLS      = "enableTLS"
	SettingCertFile       = "certFile"
	SettingKeyFile        = "keyFile"
)

var logger = log.ChildLogger(log.RootLogger(), "flow-state")

func init() {
	_ = service.RegisterFactory(&StateServiceFactory{})
}

type StateServiceFactory struct {
}

func (s *StateServiceFactory) NewService(config *service.Config) (service.Service, error) {
	ss := &StateService{}

	err := ss.init(config.Settings)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

// StateService is an implementation of StateService service
// that can access flows via URI
type StateService struct {
	server *Server
}

func (ss *StateService) Name() string {
	return "FlowStateService"
}

//DEPRECATED
func (ss *StateService) Enabled() bool {
	return true
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
func (ss *StateService) init(settings map[string]interface{}) error {

	sPort, set := settings[SettingPort]
	if !set {
		return fmt.Errorf("StateRecorder: required setting 'port' not set")
	}
	port, err := coerce.ToInt(sPort)
	if err != nil {
		return fmt.Errorf("StateRecorder: invalid port '%v'", sPort)
	}

	router := httprouter.New()

	exposeRecorder := false
	if sExpose, set := settings[SettingExposeRecorder]; set {
		exposeRecorder, _ = coerce.ToBool(sExpose)
	}

	streamingStep := false
	if stream, set := settings[SettingStreamingStep]; set {
		streamingStep, _ = coerce.ToBool(stream)
	}
	AppendEndpoints(router, logger, exposeRecorder, streamingStep)

	var options []func(*Server)

	enableTLS := false
	if sEnableTLS, set := settings[SettingEnableTLS]; set {
		enableTLS, _ = coerce.ToBool(sEnableTLS)
	}

	if enableTLS {
		certFile := ""
		if sCertFile, set := settings[SettingCertFile]; set {
			certFile, _ = coerce.ToString(sCertFile)
		}
		keyFile := ""
		if sKeyFile, set := settings[SettingKeyFile]; set {
			keyFile, _ = coerce.ToString(sKeyFile)
		}

		options = append(options, TLS(certFile, keyFile))
	}

	options = append(options, Logger(logger))

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
	})

	server, err := newServer(":"+strconv.Itoa(port), c.Handler(router), options...)
	if err != nil {
		return err
	}

	ss.server = server

	return nil
}
