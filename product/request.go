package product

import (
	"web-example/validator"
)

type Request struct {
	Name        string `validate:"required;min=3"`
	Description string
	Price       float64
	Currency    Currency `validate:"required"`
	Quantity    int      `validate:"required"`
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
