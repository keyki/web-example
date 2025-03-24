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
