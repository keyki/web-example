package order

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"web-example/product"
	"web-example/user"
	"web-example/util"
)

type Repository interface {
	ListAll(userId int) ([]*Order, error)
	Create(order *Order) error
}

type Handler struct {
	store        Repository
	userStore    user.Repository
	productStore product.Repository
}

func NewHandler(store Repository, userStore user.Repository, productStore product.Repository) *Handler {
	return &Handler{
		store:        store,
		userStore:    userStore,
		productStore: productStore,
	}
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	username, _, _ := r.BasicAuth()
	log.Println("List all orders for user: ", username)
	userInDb, _ := h.userStore.FindByUsername(username)
	orders, err := h.store.ListAll(userInDb.ID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
	log.Printf("%d orders found for user: %s", len(orders), username)
	util.WriteJSON(w, http.StatusOK, convertToResponse(orders))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err)
		util.WriteError(w, http.StatusInternalServerError, errors.New("Internal error happened"))
	}

	var req Request
	if err = json.Unmarshal(bodyBytes, &req); err != nil {
		log.Printf("Error unmarshalling request: %v", err)
		util.WriteError(w, http.StatusBadRequest, errors.New("Invalid request"))
		return
	}
	log.Printf("Received order request: %+v", req)
}

func convertToResponse(orders []*Order) []*Response {
	response := make([]*Response, 0)
	for _, o := range orders {
		response = append(response, o.ToResponse())
	}
	return response
}
