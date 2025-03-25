package order

import (
	"fmt"
	"web-example/product"
	"web-example/user"
)

type Order struct {
	ID       int                `gorm:"primaryKey,autoIncrement"`
	Products []*product.Product `gorm:"many2many:order_products"`
	UserID   int                `gorm:"not null"`
	User     *user.User         `gorm:"foreignKey:UserID"`
}

type Request struct {
	username string
	Products []*ProductRequest `json:"products"`
}

type Response struct {
	ID       int                `json:"id"`
	Products []*ProductResponse `json:"products"`
	Error    string             `json:"error"`
}

type ProductRequest struct {
	Name     string `validate:"required;min=3"`
	Quantity int    `validate:"required"`
}

type ProductResponse struct {
	Name     string
	Quantity int
}

func (r Request) String() string {
	var result string
	for _, p := range r.Products {
		result += fmt.Sprintf("\nName: %s Quantity: %d", p.Name, p.Quantity)
	}
	return result
}

func (o *Order) ToResponse() *Response {
	return &Response{
		ID:       o.ID,
		Products: convertProductsToResponses(o.Products),
	}
}

func (pr *ProductRequest) ToProductResponse() *ProductResponse {
	return &ProductResponse{
		Name:     pr.Name,
		Quantity: pr.Quantity,
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

func convertProductsToResponses(products []*product.Product) []*ProductResponse {
	response := make([]*ProductResponse, 0)
	for _, p := range products {
		response = append(response, &ProductResponse{
			Name:     p.Name,
			Quantity: p.Quantity,
		})
	}
	return response
}
