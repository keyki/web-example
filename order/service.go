package order

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"web-example/database"
	"web-example/product"
	"web-example/user"
	"web-example/util"
)

func PlaceOrder(request *Request, userStore user.Repository, productStore product.Repository,
	orderStore Repository, txService database.Transactional) (*Response, error) {

	products, err := productStore.FindAllByName(request.AllProductNames())
	if err != nil {
		log.Printf("Error finding products: %v", err)
		return nil, util.NewInternalError()
	}

	if eligibility, valid := validateOrderEligibility(products, request); !valid {
		return eligibility, nil
	}

	userInDb, err := userStore.FindByUsername(request.username)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return nil, util.NewInternalError()
	}

	order := &Order{
		Products: convertProductsToOrderProducts(products, request),
		UserID:   userInDb.ID,
	}

	tx := txService.BeginTransaction()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return nil, util.NewInternalError()
	}

	id, err := orderStore.Create(order, tx)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		rollbackTransaction(tx)
		return nil, util.NewInternalError()
	}
	log.Printf("Created order: %v", id)

	for _, prod := range products {
		prod.Quantity -= request.GetProductRequestByName(prod.Name).Quantity
		err := productStore.UpdateQuantity(prod, tx)
		if err != nil {
			log.Printf("Error updating product: %v", err)
			rollbackTransaction(tx)
			var transactionError database.TransactionError
			if errors.As(err, &transactionError) {
				return &Response{
					Error: "Product quantity has changed, please try again: " + prod.Name,
				}, nil
			}
			return nil, util.NewInternalError()
		}
	}

	result := tx.Commit()
	if result.Error != nil {
		log.Printf("Error committing transaction: %v", result.Error)
		rollbackTransaction(tx)
		return nil, util.NewInternalError()
	}

	return &Response{
		ID:       id,
		Products: convertProductRequestsToResponses(products, request.Products),
		Total:    calcTotalPriceOfRequest(products, request),
		Currency: products[0].Currency,
		Error:    "",
	}, nil
}

func rollbackTransaction(tx *gorm.DB) {
	log.Printf("Rollback transaction")
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

func validateOrderEligibility(products []*product.Product, request *Request) (*Response, bool) {
	missingProductsNames := getMissingProductNames(products, request)
	if len(missingProductsNames) > 0 {
		return &Response{
			Error: fmt.Sprintf("There are missing products: %v", missingProductsNames),
		}, false
	}

	for _, p := range products {
		if !checkProductQunatity(p, request.GetProductRequestByName(p.Name)) {
			msg := fmt.Sprintf("There are not enough product '%s' to place the order, max item(s): %d", p.Name, p.Quantity)
			log.Println(msg)
			return &Response{
				Error: msg,
			}, false
		}
	}
	return nil, true
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

func checkProductQunatity(product *product.Product, request *ProductRequest) bool {
	if product.Quantity >= request.Quantity {
		return true
	}
	return false
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
			RequestedQuantity: request.GetProductRequestByName(prd.Name).Quantity,
		})
	}
	return orderProducts
}
