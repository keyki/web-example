package order

import (
	"web-example/product"
	"web-example/user"
)

type Order struct {
	ID       int                `gorm:"primaryKey,autoIncrement"`
	Products []*product.Product `gorm:"many2many:order_products"`
	UserID   int                `gorm:"not null"`
	User     *user.User         `gorm:"foreignKey:UserID"`
}

func (o *Order) ToResponse() *Response {
	return &Response{
		ID:       o.ID,
		Products: product.ConvertToResponse(o.Products),
	}
}

type Request struct {
	Products []*ProductRequest `json:"products"`
}

type Response struct {
	ID       int                 `json:"id"`
	Products []*product.Response `json:"products"`
}

type ProductRequest struct {
	Name     string `validate:"required;min=3"`
	Quantity int    `validate:"required"`
}
