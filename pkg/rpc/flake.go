package rpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbr "github.com/prettykingking/snowflake/genproto/apis/snowflake/v1/resources"
	pbs "github.com/prettykingking/snowflake/genproto/apis/snowflake/v1/services"
	"github.com/prettykingking/snowflake/pkg/snowflake"
)

type FlakeServiceServer struct {
	pbs.UnimplementedFlakeServiceServer
}

func (p *FlakeServiceServer) GetFlake(context.Context, *pbs.GetFlakeRequest) (*pbr.Flake, error) {
	flake := pbr.Flake{}

	id := <-sCtx.flakeChan
	if id == 0 {
		return &flake, status.Errorf(codes.Unavailable, "service unavailable")
	}

	flake.IdInt = id
	flake.IdHex = snowflake.ToHex(id)

	return &flake, nil
}
