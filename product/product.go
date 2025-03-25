package product

import (
	"fmt"
	"gorm.io/plugin/optimisticlock"
	"web-example/types"
	"web-example/validator"
)

type Product struct {
	ID          int            `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"not null;index;unique"`
	Description string         `gorm:"type:text"`
	Price       float64        `gorm:"not null;default:0"`
	Currency    types.Currency `gorm:"not null;type:varchar(5);default:HUF"`
	Quantity    int            `gorm:"not null;default:1"`
	Version     optimisticlock.Version
}

type Request struct {
	Name        string `validate:"required;min=3"`
	Description string
	Price       float64
	Currency    types.Currency `validate:"required"`
	Quantity    int            `validate:"required"`
}

type Response struct {
	Name        string
	Description string
	Price       float64
	Currency    types.Currency
	Quantity    int
}

func (p *Product) ToResponse() *Response {
	return &Response{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Currency:    p.Currency,
		Quantity:    p.Quantity,
	}
}

func FindByName(products []*Product, name string) *Product {
	for _, product := range products {
		if product.Name == name {
			return product
		}
	}
	return nil
}

func (p Product) String() string {
	return fmt.Sprintf("Name: %s Price: %f %s Quantity: %d", p.Name, p.Price, p.Currency, p.Quantity)
}

func (p *Request) Validate() error {
	return validator.Validate(p)
}

func (p *Request) ToProduct() *Product {
	return &Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Currency:    p.Currency,
		Quantity:    p.Quantity,
	}
}
