package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
	"time"
	"web-example/audit"
	pb "web-example/audit/generated"
	"web-example/database"
	"web-example/log"
	"web-example/order"
	"web-example/product"
	"web-example/types"
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

	go startAuditServer(types.AuditServerPort)
	web.NewApiServer(types.WebServerPort, db).Listen()
}

func startAuditServer(port int) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.BaseLogger().Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterAuditServer(grpcServer, &audit.Server{})

	log.BaseLogger().Infof("Audit server is running on port %d", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.BaseLogger().Fatalf("Failed to serve: %v", err)
	}
}
