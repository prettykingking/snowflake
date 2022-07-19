package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/prettykingking/snowflake/genproto/apis/snowflake/v1/services"
	"github.com/prettykingking/snowflake/pkg/config"
	"github.com/prettykingking/snowflake/pkg/log"
	"github.com/prettykingking/snowflake/pkg/rpc"
	"github.com/prettykingking/snowflake/pkg/safe"
)

type Provider struct {
}

func (p *Provider) Init() error {
	return nil
}

func (p *Provider) Provide(cf *config.Configuration, rp *safe.Pool) error {
	logger, err := log.NewLogger(cf.Logging)
	if err != nil {
		return err
	}

	err = rpc.InitCtx(cf, rp)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cf.Server.Port))
	if err != nil {
		logger.Error(fmt.Sprintf("could not to listen tcp :%d", cf.Server.Port), zap.Error(err))
		return err
	}

	s := grpc.NewServer()
	// register services
	pb.RegisterFlakeServiceServer(s, &rpc.FlakeServiceServer{})

	failed := make(chan struct{})
	rp.GoCtx(func(ctx context.Context) {
		if err = s.Serve(lis); err != nil {
			logger.Error("failed to start gRPC server", zap.Error(err))
			failed <- struct{}{}
		}
	})

	rp.GoCtx(func(ctx context.Context) {
		<-ctx.Done()

		// close resources
		s.Stop()
		_ = lis.Close()

		logger.Debug("gRPC server stopped")

		_ = logger.Sync()
	})

	// timeout server start time in 2 seconds, if no failed signal received within timeout seconds
	// then assume server has started.
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			logger.Info(fmt.Sprintf("gRPC server started at %d", cf.Server.Port))
			return nil
		case <-failed:
			ticker.Stop()
			return errors.New("server failed to start")
		}
	}
}
