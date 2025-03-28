package audit

import (
	"context"
	"google.golang.org/grpc/metadata"
	pb "web-example/audit/generated"
	"web-example/log"
	"web-example/types"
)

type Server struct {
	pb.UnimplementedAuditServer
}

func (Server) LogOrder(ctx context.Context, request *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	ctx = setRequestIdIfExists(ctx)
	log.Logger(ctx).Infof("Received order to audit %s", request)
	return &pb.CreateOrderResponse{
		Result: &pb.CreateOrderResponse_Id{
			Id: request.Order.Id,
		},
	}, nil
}

func setRequestIdIfExists(ctx context.Context) context.Context {
	md, foundRequestId := metadata.FromIncomingContext(ctx)
	if foundRequestId {
		requestID := md.Get(string(types.ContextKeyReqID))
		if len(requestID) > 0 {
			ctx = context.WithValue(ctx, types.ContextKeyReqID, requestID)
		}
	}
	return ctx
}
