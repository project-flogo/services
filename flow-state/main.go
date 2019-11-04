package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/service"
	"github.com/project-flogo/services/flow-state/server/rest"
)

var port = flag.String("p", "9190", "The port of the server")

func init() {
	flag.Parse() // get the arguments from command line
}

func main() {

	logger := log.ChildLogger(log.RootLogger(), "FlowStateService")
	//for new use REST StateService
	cfg := &service.Config{
		Settings: map[string]interface{}{rest.SettingPort: *port, rest.SettingExposeRecorder:true},
	}

	logger.Info("Starting...")

	sf := &rest.StateServiceFactory{}

	s, err := sf.NewService(cfg)
	if err != nil {
		logger.Errorf("Failed to Flow State Service: %v\n", err)
		os.Exit(1)
	}

	err = s.Start()
	if err != nil {
		logger.Errorf("Failed to Flow State Service: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Started")

	exitChan := setupSignalHandling()

	code := <-exitChan

	_ = s.Stop()

	os.Exit(code)
}

func setupSignalHandling() chan int {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int, 1)
	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			exitChan <- 0
		default:
			log.RootLogger().Debugf("Unknown signal.")
			exitChan <- 1
		}
	}
	return exitChan
}
