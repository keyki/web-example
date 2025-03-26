package order

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"web-example/database"
	"web-example/product"
	"web-example/user"
	"web-example/util"
)

type Handler struct {
	store        Repository
	userStore    user.Repository
	productStore product.Repository
	txService    database.Transactional
	orderQueue   chan *CreateMessage
}

func NewHandler(store Repository, userStore user.Repository, productStore product.Repository, txService *database.TransactionService) *Handler {
	queue := make(chan *CreateMessage, 100)
	go ProcessOrder(queue, store, productStore, txService)
	return &Handler{
		orderQueue:   queue,
		store:        store,
		userStore:    userStore,
		productStore: productStore,
		txService:    txService,
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
	util.WriteJSON(w, http.StatusOK, convertOrdersToResponse(orders))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err)
		util.WriteError(w, http.StatusInternalServerError, util.NewInternalError())
	}

	var req Request
	if err = json.Unmarshal(bodyBytes, &req); err != nil {
		log.Printf("Error unmarshalling request: %v", err)
		util.WriteError(w, http.StatusBadRequest, errors.New("Invalid request"))
		return
	}
	log.Printf("Received order request: %+v", req)

	req.username = util.GetUsername(r)
	response, err := PlaceOrder(&req, h.userStore, h.productStore, h.orderQueue)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	util.WriteJSON(w, http.StatusOK, response)
}

func convertOrdersToResponse(orders []*Order) []*Response {
	response := make([]*Response, 0)
	for _, o := range orders {
		response = append(response, o.ToResponse())
	}
	return response
}
