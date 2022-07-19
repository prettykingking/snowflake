package rpc

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/prettykingking/snowflake/genproto/apis/snowflake/v1/services"
)

func TestGetFlake(t *testing.T) {
	// Set up connection to the server.
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("could not connect to gRPC server %e", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c := pb.NewFlakeServiceClient(conn)

	consumer := make(chan uint64)
	set := make(map[uint64]struct{})

	req := func() {
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		r, err := c.GetFlake(ctx, &pb.GetFlakeRequest{})
		if err != nil {
			gerr, _ := status.FromError(err)
			t.Errorf("code=%d message=%s", gerr.Code(), gerr.Message())

			consumer <- 0
		} else {
			consumer <- r.IdInt
		}
	}

	for i := 0; i < 10000; i++ {
		go req()
	}

	for i := 0; i < 10000; i++ {
		id := <-consumer
		if id == 0 {
			t.Errorf("service response with error")
			break
		}

		if _, ok := set[id]; ok {
			t.Fatal("duplicated id")
		}

		set[id] = struct{}{}
	}
}
