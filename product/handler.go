package product

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"web-example/util"
)

type Handler struct {
	store Repository
}

func NewHandler(store Repository) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.ListAll()
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
	util.WriteJSON(w, http.StatusOK, ConvertToResponse(products))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var productRequest Request
	if err := json.NewDecoder(r.Body).Decode(&productRequest); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := productRequest.Validate(); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.store.Create(productRequest.ToProduct()); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		log.Printf("Create Error: %v", err)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	log.Printf("Find product %s\n", name)
	if name == "" {
		util.WriteError(w, http.StatusBadRequest, errors.New("Name is required"))
		return
	}

	product, err := h.store.FindByName(name)
	if err != nil {
		log.Printf("Find error: %v", err)
		util.WriteJSON(w, http.StatusNotFound, []*Response{})
		return
	}

	util.WriteJSON(w, http.StatusOK, product.ToResponse())
}

func ConvertToResponse(products []*Product) (r []*Response) {
	for _, p := range products {
		r = append(r, p.ToResponse())
	}
	return r
}
