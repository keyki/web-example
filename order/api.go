package order

import (
	"fmt"
	"web-example/product"
	"web-example/types"
)

type Request struct {
	username string
	Products []*ProductRequest `json:"products"`
}

type Response struct {
	ID       int                `json:"id"`
	Products []*ProductResponse `json:"products"`
	Total    float64            `json:"total"`
	Currency types.Currency     `json:"currency"`
	Error    string             `json:"error"`
}

type ProductRequest struct {
	Name     string `validate:"required;min=3"`
	Quantity int    `validate:"required"`
}

type ProductResponse struct {
	Name     string
	Quantity int
	Price    float64
	Currency types.Currency
}

func (r Request) String() string {
	var result string
	for _, p := range r.Products {
		result += fmt.Sprintf("Name: %s Quantity: %d ", p.Name, p.Quantity)
	}
	return result
}

func (pr *ProductRequest) ToProductResponse(prod *product.Product) *ProductResponse {
	return &ProductResponse{
		Name:     pr.Name,
		Quantity: pr.Quantity,
		Price:    prod.Price,
		Currency: prod.Currency,
	}
}

func (r *Request) AllProductNames() []string {
	names := make([]string, 0)
	for _, p := range r.Products {
		names = append(names, p.Name)
	}
	return names
}

func (r *Request) GetProductRequestByName(name string) *ProductRequest {
	for _, p := range r.Products {
		if p.Name == name {
			return p
		}
	}
	return nil
}
