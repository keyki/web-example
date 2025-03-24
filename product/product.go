package product

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
    Quantity    int      `json:"quantity"`
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
