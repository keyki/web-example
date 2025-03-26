package product

import (
	"fmt"
	"gorm.io/plugin/optimisticlock"
	"web-example/types"
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
