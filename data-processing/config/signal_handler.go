package config

import (
	"os"
	"os/signal"
	"syscall"
)

func InitSignalHandler() chan bool {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	doneChan := make(chan bool, 1)

	go func() {
		sig := <-sigChan
		log.Info("Received signal: ", sig)
		doneChan <- true
		log.Info("Exiting...?")
	}()

	return doneChan
}