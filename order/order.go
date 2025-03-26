package order

import (
	"fmt"
	"web-example/product"
	"web-example/user"
)

type Order struct {
	ID       int             `gorm:"primaryKey,autoIncrement"`
	Products []*OrderProduct `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	UserID   int             `gorm:"not null"`
	User     *user.User      `gorm:"foreignKey:UserID"`
}

type OrderProduct struct {
	OrderID           int              `gorm:"primaryKey"`
	ProductID         int              `gorm:"primaryKey"`
	ProductName       string           `gorm:"not null"`
	Product           *product.Product `gorm:"foreignKey:ProductID"`
	RequestedQuantity int              `gorm:"not null"`
}

func (o Order) String() string {
	return fmt.Sprintf("Order: products: %v", o.Products)
}

func (o OrderProduct) String() string {
	return fmt.Sprintf("Product: name: %s, requested: %d", o.ProductName, o.RequestedQuantity)
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
		if p.ProductName == name {
			return p.Product
		}
	}
	return nil
}

func (o *Order) FindOrderProductByName(name string) *OrderProduct {
	for _, p := range o.Products {
		if p.ProductName == name {
			return p
		}
	}
	return nil
}

func (o *Order) AllProductIds() []int {
	ids := make([]int, 0)
	for _, p := range o.Products {
		ids = append(ids, p.ProductID)
	}
	return ids
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
