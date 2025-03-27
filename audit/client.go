package audit

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"time"
	pb "web-example/audit/generated"
	"web-example/log"
)

type Client struct {
	grpcClient pb.AuditClient
}

func NewClient() *Client {
	logger := log.BaseLogger()
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Failed to create gRPC client: %v", err)
	}

	for {
		state := conn.GetState()
		logger.Info("gRPC audit client connection state:", state)

		if state == connectivity.Ready {
			break
		}

		conn.WaitForStateChange(context.Background(), state)
		time.Sleep(time.Second)
	}

	client := pb.NewAuditClient(conn)
	return &Client{
		grpcClient: client,
	}
}

func (c *Client) LogOrder(ctx context.Context, message *pb.CreateOrderRequest) {
	logger := log.Logger(ctx)
	logger.Infof("Sending order to audit: %v", message)

	response, err := c.grpcClient.LogOrder(ctx, message)
	if err != nil {
		logger.Infof("Failed to send order: %v", err)
		return
	}

	if len(response.GetError()) != 0 {
		logger.Infof("Failed to log order: %v", response.GetError())
		return
	}

	logger.Infof("Successfully audited order: %v", response.GetId())
}
