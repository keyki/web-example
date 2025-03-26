package order

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"web-example/database"
	"web-example/log"
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
	logger := log.Logger(r.Context())
	logger.Info("List all orders for user: ", username)
	userInDb, _ := h.userStore.FindByUsername(r.Context(), username)
	orders, err := h.store.ListAll(r.Context(), userInDb.ID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
	logger.Infof("%d orders found for user: %s", len(orders), username)
	util.WriteJSON(w, http.StatusOK, convertOrdersToResponse(orders))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Logger(r.Context()).Infof("Error reading body: %s", err)
		util.WriteError(w, http.StatusInternalServerError, util.NewInternalError())
	}

	var req Request
	if err = json.Unmarshal(bodyBytes, &req); err != nil {
		log.Logger(r.Context()).Infof("Error unmarshalling request: %v", err)
		util.WriteError(w, http.StatusBadRequest, errors.New("Invalid request"))
		return
	}
	log.Logger(r.Context()).Infof("Received order request: %+v", req)

	req.username = util.GetUsername(r)
	response, err := PlaceOrder(r.Context(), &req, h.userStore, h.productStore, h.orderQueue)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	util.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	logger := log.Logger(r.Context())
	orderId := r.PathValue("orderId")
	logger.Infof("Received order delete request: %+v", orderId)

	id, err := strconv.Atoi(orderId)
	if err != nil {
		logger.Infof("Error converting order id to int: %v", err)
		util.WriteError(w, http.StatusBadRequest, errors.New("Invalid orderId"))
		return
	}

	userInDb, err := h.userStore.FindByUsername(r.Context(), util.GetUsername(r))
	if err != nil {
		logger.Infof("Error finding user: %v", err)
		util.WriteError(w, http.StatusInternalServerError, err)
	}

	order, err := h.store.Find(r.Context(), id, userInDb.ID)
	if err != nil {
		logger.Infof("Error finding order: %v", err)
		util.WriteError(w, http.StatusInternalServerError, util.NewInternalError())
		return
	}

	if order == nil {
		logger.Infof("Order with id %d not found", id)
		util.WriteError(w, http.StatusNotFound, errors.New("Order not found"))
		return
	}

	err = h.store.Delete(r.Context(), order)
	if err != nil {
		logger.Infof("Error deleting order: %v", err)
		util.WriteError(w, http.StatusInternalServerError, util.NewInternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func convertOrdersToResponse(orders []*Order) []*Response {
	response := make([]*Response, 0)
	for _, o := range orders {
		response = append(response, o.ToResponse())
	}
	return response
}
