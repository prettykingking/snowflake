package server

import (
	"context"
	"os/signal"
	"syscall"
)

func (s *Server) configureSignals() {
	signal.Notify(s.signalChan, syscall.SIGUSR1)
}

// listenSignals listens for POSIX signals
func (s *Server) listenSignals(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-s.signalChan:
			if sig == syscall.SIGUSR1 {
				// this signal usually used to Rotate log
			}
		}
	}
}
