package rpc

import (
	"context"
	"errors"
	"runtime"
	"time"

	"go.uber.org/zap"

	"github.com/prettykingking/snowflake/pkg/config"
	"github.com/prettykingking/snowflake/pkg/log"
	"github.com/prettykingking/snowflake/pkg/safe"
	"github.com/prettykingking/snowflake/pkg/snowflake"
)

type serviceContext struct {
	snowflake *snowflake.Snowflake
	flakeChan chan uint64
}

var sCtx = serviceContext{
	flakeChan: make(chan uint64),
}

// InitCtx resources necessary for services
func InitCtx(cf *config.Configuration, rp *safe.Pool) error {
	err := initSnowflake(cf.Settings, rp)
	if err != nil {
		return err
	}

	return nil
}

func initSnowflake(cf *config.Settings, rp *safe.Pool) error {
	st := snowflake.Settings{}
	st.MachineId = func() (uint16, error) {
		return cf.MachineId, nil
	}

	start, err := time.Parse("2006-01-02T15:04:05Z", cf.StartTime)
	if err != nil {
		return err
	}

	st.StartTime = start

	sf := snowflake.NewSnowflake(st)
	if sf == nil {
		return errors.New("snowflake not created")
	}

	logger := log.GetLogger()
	logger.Info("snowflake settings", zap.String("startTime", cf.StartTime),
		zap.Uint16("machineId", cf.MachineId))

	sCtx.snowflake = sf

	for i := 0; i <= runtime.NumCPU(); i++ {
		rp.GoCtx(func(ctx context.Context) {
			id, err := sCtx.snowflake.NextID()
			if err != nil {
				return
			}

			for {
				select {
				case <-ctx.Done():
					return
				case sCtx.flakeChan <- id:
					id, err = sCtx.snowflake.NextID()
					if err != nil {
						return
					}
				}
			}
		})
	}

	return nil
}
