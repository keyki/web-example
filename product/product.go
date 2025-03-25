package product

import "web-example/validator"

var (
	HUF = Currency("HUF")
	USD = Currency("USD")
	EUR = Currency("EUR")
)

type Currency string

type Product struct {
	ID          int      `gorm:"primaryKey;autoIncrement"`
	Name        string   `gorm:"not null;index;unique"`
	Description string   `gorm:"type:text"`
	Price       float64  `gorm:"not null;default:0"`
	Currency    Currency `gorm:"not null;type:varchar(5);default:HUF"`
	Quantity    int      `gorm:"not null;default:1"`
}

type Request struct {
	Name        string `validate:"required;min=3"`
	Description string
	Price       float64
	Currency    Currency `validate:"required"`
	Quantity    int      `validate:"required"`
}

type Response struct {
	Name        string
	Description string
	Price       float64
	Currency    Currency
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
