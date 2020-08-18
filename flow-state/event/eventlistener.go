package event

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/project-flogo/core/engine/event"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/flow/state"
)

var recorderLog = log.ChildLogger(log.RootLogger(), "step-listener")

const (
	// PongWait Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// PingPeriod Send pings to peer with this period. Must be less than PongWait.
	pingPeriod = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var stepEventQueue = make(chan *event.Context, 10)
var stepChan = make(chan *state.Step)

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
				stepChan <- t.step
			}
		}
	}
}

// Status is a basic health check for the server to determine if it is up
func HandleStepEvent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	recorderLog.Debugf("Received step event websocket request: %+v", r)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		recorderLog.Errorf("websocket upgrade failed: %s", err.Error())
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			recorderLog.Errorf("close websocket error: %v", err)
		}
	}()
	// Send ping
	ticker := time.NewTicker(pingPeriod)
	conn.SetPingHandler(func(string) error {
		if recorderLog.DebugEnabled() {
			recorderLog.Debug("Ping Handler")
		}
		if err := conn.WriteControl(websocket.PongMessage, nil, time.Now().Add(pongWait)); err != nil {
			recorderLog.Warnf("Sending PONG to client error: %v", err)
		}
		return nil
	})

	conn.SetPongHandler(func(string) error {
		if recorderLog.DebugEnabled() {
			recorderLog.Debug("PONG Handler")
		}
		conn.SetReadDeadline(time.Now().Add(pongWait)) // reset read timeout on receiving a pong
		return nil
	})

	for {
		select {
		case step := <-stepChan:
			err = conn.WriteJSON(step)
			if err != nil {
				conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInternalServerErr, fmt.Sprintf("Write json data error:%s", err.Error())),
					time.Now())
				recorderLog.Error("error writing message: %s", err.Error())
				return
			}
		case <-ticker.C:
			recorderLog.Debugf("Sending heartbeat Ping to client")
			// NOTE: Control frame writes do not need to be synchronized
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(pongWait)); err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					recorderLog.Warnf("Sent heartbeat Ping timeout: %v", err)
				} else {
					// for any other error we need to close tunnel
					recorderLog.Warnf("Sent heartbeat Ping error: %v", err)
					conn.Close()
				}
			}
		}
	}
}
