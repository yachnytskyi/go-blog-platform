package application

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	shutDownOnce sync.Once
	terminate    = make(chan os.Signal, 1)
	waitGroup    sync.WaitGroup
)

func GracefulShutdown() {
	shutDownOnce.Do(func() {
		signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
		waitGroup.Add(1)

		go func() {
			<-terminate
			os.Exit(1)
		}()
	})
	waitGroup.Wait()
}
