package server

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/prettykingking/snowflake/pkg/safe"
)

type Server struct {
	signalChan   chan os.Signal
	stopChan     chan bool
	routinesPool *safe.Pool
}

// NewServer returns new server instance
func NewServer(routinesPool *safe.Pool) *Server {
	svr := &Server{
		routinesPool: routinesPool,
		signalChan:   make(chan os.Signal, 1),
		stopChan:     make(chan bool, 1),
	}

	svr.configureSignals()

	return svr
}

// Start starts application server
func (s *Server) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	s.routinesPool.GoCtx(s.listenSignals)
}

// Wait blocks until the server shutdown
func (s *Server) Wait() {
	<-s.stopChan
}

// Stop stops the server
func (s *Server) Stop() {
	s.stopChan <- true
}

// Close destroys the server resources
func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	go func(ctx context.Context) {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			panic("Timeout while stopping server, killing instance ✝")
		}
	}(ctx)

	s.routinesPool.Stop()

	signal.Stop(s.signalChan)

	close(s.signalChan)
	close(s.stopChan)

	cancel()
}
