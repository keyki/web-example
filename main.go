package main

import (
	"time"
	"web-example/database"
	"web-example/log"
	"web-example/order"
	"web-example/product"
	"web-example/user"
	"web-example/web"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5433 sslmode=disable"
	db, err := database.Connect(dsn, &database.Options{
		MaxOpenConns:    10,
		MaxIdleConns:    10,
		ConnMaxLifetime: 10 * time.Minute,
	})
	if err != nil {
		log.BaseLogger().Fatalf("Failed to connect to database: %v", err)
	}
	log.BaseLogger().Println("Successfully connected to database")

	err = db.AutoMigrate(
		&user.User{},
		&product.Product{},
		&order.Order{},
		&order.OrderProduct{},
	)
	if err != nil {
		log.BaseLogger().Fatalf("Failed to migrate schema: %v", err)
	}

	server := web.NewApiServer(8080, db)
	server.Listen()
}
