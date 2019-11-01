package local

import (
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/service"
	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/store"
)

func init() {
	_ = service.RegisterFactory(&StateRecorderFactory{})
}

type StateRecorderFactory struct {
}

func (s StateRecorderFactory) NewService(config *service.Config) (service.Service, error) {
	recorder := &StateRecorder{}
	recorder.init(config.Settings)

	//todo switch this logger
	recorder.logger = log.RootLogger()

	return recorder, nil
}

// StateRecorder is an implementation of StateRecorder service
// that can access flows via URI
type StateRecorder struct {
	stepStore store.StepStore
	snapshotStore store.SnapshotStore
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
func (sr *StateRecorder) init(settings map[string]interface{}) {
	sr.snapshotStore = store.GetSnapshotStore()
	sr.stepStore = store.GetStepStore()
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
	if sr.stepStore != nil {
		err := sr.stepStore.SaveStep(step)
		if err != nil {
			return err
		}
	}
	return nil
}
