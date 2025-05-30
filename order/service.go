package order

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"web-example/audit"
	pb "web-example/audit/generated"
	"web-example/log"
	"web-example/product"
	"web-example/user"
	"web-example/util"
)

func PlaceOrder(ctx context.Context, request *Request, userStore user.Repository,
	productStore product.Repository, queue chan *CreateMessage, auditClient *audit.Client) (*Response, error) {
	logger := log.Logger(ctx)

	products, err := productStore.FindAllByName(ctx, request.AllProductNames())
	if err != nil {
		logger.Infof("Error finding products: %v", err)
		return nil, util.NewInternalError()
	}

	if len(products) != len(request.Products) {
		missingProductNames := getMissingProductNames(products, request)
		return &Response{
			Error: fmt.Sprintf("There are missing products: %v", missingProductNames),
		}, nil
	}

	userInDb, err := userStore.FindByUsername(ctx, request.username)
	if err != nil {
		logger.Infof("Error finding user: %v", err)
		return nil, util.NewInternalError()
	}

	order := &Order{
		Products: convertProductsToOrderProducts(products, request),
		UserID:   userInDb.ID,
	}

	errResp := make(chan error)
	idResp := make(chan int)
	queue <- &CreateMessage{
		Order:       order,
		ErrResponse: errResp,
		IdResponse:  idResp,
		Context:     ctx,
	}

	select {
	case err := <-errResp:
		logger.Infof("Error creating order: %v", err)
		return &Response{
			Error: err.Error(),
		}, nil
	case idResp := <-idResp:
		auditClient.LogOrder(ctx, &pb.CreateOrderRequest{
			Order: &pb.Order{
				Id:       int32(idResp),
				UserId:   int32(userInDb.ID),
				Products: convertToAuditProducts(order.Products),
			},
		})
		return &Response{
			ID:       idResp,
			Products: convertProductRequestsToResponses(products, request.Products),
			Total:    calcTotalPriceOfRequest(products, request),
			Currency: products[0].Currency,
			Error:    "",
		}, nil
	}
}

func rollbackTransaction(ctx context.Context, tx *gorm.DB) {
	log.Logger(ctx).Info("Rollback transaction")
	tx.Rollback()
}

func calcTotalPriceOfRequest(products []*product.Product, request *Request) float64 {
	price := 0.0
	for _, prd := range products {
		prodReq := request.GetProductRequestByName(prd.Name)
		price += prd.Price * float64(prodReq.Quantity)
	}
	return price
}

func convertProductRequestsToResponses(products []*product.Product, requests []*ProductRequest) []*ProductResponse {
	response := make([]*ProductResponse, 0)
	for _, request := range requests {
		prod := product.FindByName(products, request.Name)
		response = append(response, request.ToProductResponse(prod))
	}
	return response
}

func convertProductsToOrderProducts(products []*product.Product, request *Request) []*OrderProduct {
	var orderProducts []*OrderProduct
	for _, prd := range products {
		orderProducts = append(orderProducts, &OrderProduct{
			ProductID:         prd.ID,
			ProductName:       prd.Name,
			RequestedQuantity: request.GetProductRequestByName(prd.Name).Quantity,
		})
	}
	return orderProducts
}

func convertToAuditProducts(products []*OrderProduct) []*pb.OrderProduct {
	var auditProducts []*pb.OrderProduct
	for _, prd := range products {
		auditProducts = append(auditProducts, &pb.OrderProduct{
			ProductId: int32(prd.ProductID),
			Quantity:  int32(prd.RequestedQuantity),
		})
	}
	return auditProducts
}

func getMissingProductNames(products []*product.Product, request *Request) []string {
	var missingProducts []string
	if len(products) != len(request.Products) {
		for _, pr := range request.Products {
			found := false
			for _, p := range products {
				if pr.Name == p.Name {
					found = true
				}
			}
			if !found {
				missingProducts = append(missingProducts, pr.Name)
			}
		}
	}
	return missingProducts
}
