package web

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"web-example/database"
	"web-example/log"
	"web-example/order"
	"web-example/product"
	"web-example/user"
)

type Server struct {
	port int
	db   *gorm.DB
}

func NewApiServer(port int, db *gorm.DB) *Server {
	return &Server{port: port, db: db}
}

func (s *Server) Listen() {
	log.BaseLogger().Infof("Listening on port %d", s.port)

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)

	transactionService := database.NewDbTransaction(s.db)
	orderStore := order.NewStore(s.db)
	orderHandler := order.NewHandler(orderStore, userStore, productStore, transactionService)

	v1Mux := http.NewServeMux()
	v1Mux.HandleFunc("GET /users", userHandler.ListAll)
	v1Mux.HandleFunc("POST /user", userHandler.Create)
	v1Mux.HandleFunc("GET /user/{userName}", userHandler.Get)

	v1Mux.HandleFunc("GET /products", productHandler.ListAll)
	v1Mux.HandleFunc("POST /product", productHandler.Create)
	v1Mux.HandleFunc("GET /product/{name}", productHandler.Get)

	v1Mux.HandleFunc("GET /orders", orderHandler.ListAll)
	v1Mux.HandleFunc("POST /order", orderHandler.Create)
	v1Mux.HandleFunc("DELETE /order/{orderId}", orderHandler.Delete)

	authMiddleware := CreateMiddleware(AuthenticationMiddleware)
	wrappedAuthMux := authMiddleware(userStore, http.StripPrefix("/api/v1", v1Mux))

	mainMux := http.NewServeMux()
	mainMux.HandleFunc("GET /info", Info)
	mainMux.Handle("/api/v1/", wrappedAuthMux)

	mainMiddleware := CreateMiddleware(RequestIdMiddleware, MeasureMiddleware)
	wrappedMainMux := mainMiddleware(nil, mainMux)

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), wrappedMainMux)
	if err != nil {
		log.BaseLogger().Fatalf("Failed to listen on port %d: %v", s.port, err)
	}
}

func Info(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"version": "v1", "status": "ok"}`))
}
