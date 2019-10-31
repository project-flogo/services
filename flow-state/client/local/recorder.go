package local

import (
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/flow/service"
	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/store"
)

// StateRecorder is an implementation of StateRecorder service
// that can access flows via URI
type StateRecorder struct {
	enabled bool
	stepStore store.StepStore
	snapshotStore store.SnapshotStore
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

	if sr.enabled {
		sr.snapshotStore = store.GetSnapshotStore()
		sr.stepStore = store.GetStepStore()
	}
}

// RecordSnapshot implements instance.StateRecorder.RecordSnapshot
func (sr *StateRecorder) RecordSnapshot(snapshot *state.Snapshot) error {

	if sr.snapshotStore != nil {
		err := sr.snapshotStore.SaveSnapshot(snapshot)
		if err != nil {
			return err
		}
	}
	return nil
}

// RecordStep implements instance.StateRecorder.RecordStep
func (sr *StateRecorder) RecordStep(step *state.Step) error {

	return nil
}


func DefaultConfig() *support.ServiceConfig {
	return &support.ServiceConfig{Name: service.ServiceStateRecorder, Enabled: true, Settings: map[string]string{"host": ""}}
}
