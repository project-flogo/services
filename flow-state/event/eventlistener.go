package event

import (
	"net/http"
	"time"


	"github.com/julienschmidt/httprouter"
	"github.com/project-flogo/core/engine/event"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/flow/state"
	"github.com/project-flogo/services/flow-state/stream"
)

var recorderLog = log.ChildLogger(log.RootLogger(), "step-listener")

var stepEventQueue = make(chan *event.Context, 1)
var done = make(chan bool, 1)
var stepChan = make(chan *state.Step, 100)

var s stream.EventStreamScheduler

type recorderEvent struct {
}

func (ls *recorderEvent) HandleEvent(evt *event.Context) error {
	stepEventQueue <- evt
	return nil
}

func StartStepListener() {
	err := event.RegisterListener("state-recorder-step-listener", &recorderEvent{}, []string{EventType})
	if err != nil {
		recorderLog.Errorf("Failed to enable state-recorder-step-listener due to error - '%v'", err)
	}
	go handleRecordEvent()
}

func handleRecordEvent() {
	for {
		select {
		case stepE := <-stepEventQueue:
			switch t := stepE.GetEvent().(type) {
			case *stepEvent:
				if s == nil {
					for checking := true; checking; {
						if s != nil {
							checking = false
							break
						}
						recorderLog.Infof("Waiting for client step steaming connection")
						time.Sleep(1 * time.Second)
					}
					recorderLog.Infof("Client steaming steps connected")
				}
				step := t.Step()
				if len(step.FlowChanges) > 0 {
					//main flow
					flowChange, ok := step.FlowChanges[0]
					if ok && (flowChange.Status == 500 || flowChange.Status == 600 || flowChange.Status == 700) {
						//Finish it after flow completed/canceled/failed
						s.Finish(&stream.Response{Step: step, Status: stream.Completed})
						s = nil
						done <- true
					} else {
						s.UpdateResponse(&stream.Response{Step: step, Status: stream.Running})
					}
				} else {
					s.UpdateResponse(&stream.Response{Step: step, Status: stream.Running})
				}
			}
		}
	}
}

func setWriter(ss stream.EventStreamScheduler) {
	s = ss
}

// Status is a basic health check for the server to determine if it is up
func HandleStepEvent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	connectionMonitor := stream.NewConnectionMonitor(make(chan int, 10))
	if cn, ok := w.(http.CloseNotifier); ok {
		// Listen to the closing of the http connection via the CloseNotifier
		notify := cn.CloseNotify()
		go func() {
			<-notify
			recorderLog.Info("Connection closed by the client or broken.")
			connectionMonitor.ConnectionLost()
		}()
	}

	scheduler := stream.NewEventStreamScheduler(w, r, &stream.Response{Status: stream.Running, Details: "Waiting step data", Step: nil})
	go func() {
		// Start the scheduler
		recorderLog.Debug("Starting connection monitor scheduler...")
		scheduler.Start(connectionMonitor)
	}()

	recorderLog.Info("Step streaming request comes")

	setWriter(scheduler)
	<-done
}
