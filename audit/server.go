package audit

import (
	"context"
	pb "web-example/audit/generated"
	"web-example/log"
)

type Server struct {
	pb.UnimplementedAuditServer
}

func (Server) LogOrder(ctx context.Context, request *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log.BaseLogger().Infof("Received order to audit %s", request)
	return &pb.CreateOrderResponse{
		Result: &pb.CreateOrderResponse_Id{
			Id: request.Order.Id,
		},
	}, nil
}
