package order

import (
	"fmt"
	"web-example/product"
	"web-example/types"
	"web-example/user"
)

type Order struct {
	ID       int             `gorm:"primaryKey,autoIncrement"`
	Products []*OrderProduct `gorm:"foreignKey:OrderID"`
	UserID   int             `gorm:"not null"`
	User     *user.User      `gorm:"foreignKey:UserID"`
}

type OrderProduct struct {
	OrderID           int              `gorm:"primaryKey"`
	ProductID         int              `gorm:"primaryKey"`
	Product           *product.Product `gorm:"foreignKey:ProductID"`
	RequestedQuantity int              `gorm:"not null"`
}

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
		result += fmt.Sprintf("\nName: %s Quantity: %d", p.Name, p.Quantity)
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

func (o Order) String() string {
	return fmt.Sprintf("Order: products: %v", o.Products)
}

func (o *Order) ToResponse() *Response {
	return &Response{
		ID:       o.ID,
		Products: convertProductsToResponses(o.Products),
		Total:    calcTotalPriceOfOrderProducts(o.Products),
		Currency: o.Products[0].Product.Currency,
	}
}

func (o *Order) FindProductByName(name string) *product.Product {
	for _, p := range o.Products {
		if p.Product.Name == name {
			return p.Product
		}
	}
	return nil
}

func calcTotalPriceOfOrderProducts(products []*OrderProduct) float64 {
	total := 0.0
	for _, prd := range products {
		total += prd.Product.Price * float64(prd.RequestedQuantity)
	}
	return total
}

func convertProductsToResponses(orderProducts []*OrderProduct) []*ProductResponse {
	response := make([]*ProductResponse, 0)
	for _, op := range orderProducts {
		response = append(response, &ProductResponse{
			Name:     op.Product.Name,
			Quantity: op.RequestedQuantity,
			Price:    op.Product.Price,
			Currency: op.Product.Currency,
		})
	}
	return response
}
