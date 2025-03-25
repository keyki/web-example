package order

import (
	"fmt"
	"log"
	"web-example/product"
	"web-example/user"
	"web-example/util"
)

func PlaceOrder(request *Request, userStore user.Repository, productStore product.Repository,
	orderStore Repository) (*Response, error) {

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
		Products: products,
		UserID:   userInDb.ID,
		User:     userInDb,
	}

	id, err := orderStore.Create(order)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return nil, util.NewInternalError()
	}
	log.Printf("Created order: %v", id)

	return &Response{
		ID:       id,
		Products: convertProductRequestsToResponses(request.Products),
		Error:    "",
	}, nil
}

func validateOrderEligibility(products []*product.Product, request *Request) (*Response, bool) {
	missingProductsNames := getMissingProductNames(products, request)
	if len(missingProductsNames) > 0 {
		return &Response{
			Error: fmt.Sprintf("There are missing products: %v", missingProductsNames),
		}, false
	}

	for _, p := range products {
		if !checkOrderEligibility(p, request.GetProductRequestByName(p.Name)) {
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

func checkOrderEligibility(product *product.Product, request *ProductRequest) bool {
	if product.Quantity >= request.Quantity {
		return true
	}
	return false
}

func convertProductRequestsToResponses(requests []*ProductRequest) []*ProductResponse {
	response := make([]*ProductResponse, 0)
	for _, request := range requests {
		response = append(response, request.ToProductResponse())
	}
	return response
}
