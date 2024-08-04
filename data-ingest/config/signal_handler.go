package config

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func InitSignalHandler() chan bool {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	doneChan := make(chan bool, 1)

	go func() {
		sig := <-sigChan
		log.Info("Received signal: ", sig)
		doneChan <- true
	}()

	return doneChan
}