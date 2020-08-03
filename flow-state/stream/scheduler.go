package stream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/project-flogo/flow/state"
	log "github.com/sirupsen/logrus"
)

//Response status
const (
	Failed    = "Failed"
	Running   = "Running"
	Completed = "Completed"
)

// Schedules the Server Side Event responses
type EventStreamScheduler interface {
	Start(ConnectionMonitor)
	UpdateResponse(res *Response)
	FinishWithError(errorDetails string)
	Finish(res *Response)
}

// Schedules the event stream responses
type eventStreamScheduler struct {
	W                        http.ResponseWriter
	R                        *http.Request
	defaultResponse          *Response
	eventStreamChannel       chan *Response
	schedulerFinishedChannel chan int
	responseInterval         int
	responseChannel          chan *Response
	running                  bool
	connectionLost           bool
}

type Response struct {
	Status  string      `json:"status"`
	Details string      `json:"details"`
	Step    *state.Step `json:"step"`
}

// NewEventStreamScheduler create new Event Stream Scheduler
func NewEventStreamScheduler(w http.ResponseWriter, r *http.Request, defResponse *Response) EventStreamScheduler {
	return &eventStreamScheduler{
		W:                        w,
		R:                        r,
		defaultResponse:          defResponse,
		eventStreamChannel:       make(chan *Response),
		schedulerFinishedChannel: make(chan int),
		responseInterval:         5000,
		responseChannel:          make(chan *Response),
		running:                  true,
		connectionLost:           false,
	}
}

// Start indicate to start the Event Stream Scheduler
func (s *eventStreamScheduler) Start(connectionMonitor ConnectionMonitor) {
	// Mark the end of the scheduler
	defer func() {
		defer func() {
			if !s.connectionLost {
				s.schedulerFinishedChannel <- 1
			}
			s.running = false
		}()
	}()

	log.Info("Event Scheduler Started...")

	s.W.Header().Set("Content-Type", "application/json")
	s.W.WriteHeader(http.StatusOK)

	fmt.Fprintf(s.W, "%s", "[")
	// Write to the ResponseWriter, `w`.
	json.NewEncoder(s.W).Encode(s.defaultResponse)

	// Flush the response.  This is only possible if
	// the response supports streaming.
	if f, ok := s.W.(http.Flusher); ok {
		log.WithFields(log.Fields{"response": s.defaultResponse}).Info("Sending Response: ")
		f.Flush()
	}

	for building := true; building; {

		select {
		case res := <-s.eventStreamChannel:
			if !connectionMonitor.IsAlive() {
				// No connection so we should not write in response
				log.Debug("Connection closed, stopping Event Stream Scheduler")
				// Finishing build
				building = false
				break
			}

			// Print separator of json response
			fmt.Fprintf(s.W, "%s", ",")
			// Write to the ResponseWriter, `w`.
			json.NewEncoder(s.W).Encode(res)
			// Print end of json response
			fmt.Fprintf(s.W, "%s", "]")

			// Flush the response.  This is only possible if
			// the response supports streaming.
			if f, ok := s.W.(http.Flusher); ok {
				log.WithFields(log.Fields{"response": res}).Debugf("Sending Response: ")
				f.Flush()
			}
			building = false
		case <-time.After(time.Millisecond * time.Duration(s.responseInterval)):
			if !connectionMonitor.IsAlive() {
				// No connection so we should not write in response
				log.Debug("Connection closed, waiting for scheduler to stop")
				// Finishing build
				building = false
				s.connectionLost = true
				continue
			}

			// Print separator of json response
			fmt.Fprintf(s.W, "%s", ",")
			// Write to the ResponseWriter, `w`.
			json.NewEncoder(s.W).Encode(s.defaultResponse)
			// Flush the response.  This is only possible if
			// the response supports streaming.
			if f, ok := s.W.(http.Flusher); ok {
				log.WithFields(log.Fields{
					"response": s.defaultResponse,
				}).Info("Sending Response: ")
				f.Flush()
			}

		case res := <-s.responseChannel:
			if !connectionMonitor.IsAlive() {
				// No connection so we should not write in response
				log.Debug("Connection closed, waiting for scheduler to stop")
				// Finishing build
				building = false
				continue
			}

			// Print separator of json response
			fmt.Fprintf(s.W, "%s", ",")
			// Write to the ResponseWriter, `w`.
			json.NewEncoder(s.W).Encode(res)
			// Flush the response.  This is only possible if
			// the response supports streaming.
			if f, ok := s.W.(http.Flusher); ok {
				log.WithFields(log.Fields{"response": s.defaultResponse}).Debugf("Sending Response: ")
				f.Flush()
			}
		}
	}

}

// UpdateResponse UpdateResponse indicate to update the response
func (s *eventStreamScheduler) UpdateResponse(res *Response) {
	s.responseChannel <- res
}

// FinishWithError indicate to complete the response with errors
func (s *eventStreamScheduler) FinishWithError(errorDetails string) {
	// Finish scheduler
	s.Finish(&Response{
		Status:  Failed,
		Details: errorDetails,
		Step:    nil,
	})
}

// Finish indicate to complete the response
func (s *eventStreamScheduler) Finish(res *Response) {
	// Send response to the stream channel
	if s.running {
		s.eventStreamChannel <- res
	}
	// Wait for the scheduler to finish
	if !s.connectionLost {
		<-s.schedulerFinishedChannel
	}
}
