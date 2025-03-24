package main

import (
    "log"
    "time"
    "web-example/database"
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
        log.Fatalf("Failed to connect to database: %v", err)
    }
    log.Println("Successfully connected to database")

    err = db.AutoMigrate(&user.User{})
    if err != nil {
        log.Fatalf("Failed to migrate user: %v", err)
    }
    log.Println("Successfully migrated users table")

    server := web.NewApiServer(8080, db)
    server.Listen()
}
