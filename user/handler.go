package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"web-example/log"
	"web-example/util"
)

type Handler struct {
	store Repository
}

func NewHandler(store Repository) *Handler {
	return &Handler{store: store}
}

func convertToUserResponse(users []*User) []*Response {
	r := make([]*Response, 0)
	for _, u := range users {
		r = append(r, u.ToResponse())
	}
	return r
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.ListAll(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
	util.WriteJSON(w, http.StatusOK, convertToUserResponse(users))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var userRequest Request
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := userRequest.Validate(); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userRequest.Password = util.HashPassword(userRequest.Password)
	if err := h.store.Create(r.Context(), userRequest.ToUser()); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		log.Logger(r.Context()).Infof("Create Error: %v", err)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("userName")
	log.Logger(r.Context()).Infof("Find user %s\n", userName)
	if userName == "" {
		util.WriteError(w, http.StatusBadRequest, errors.New("UserName is required"))
		return
	}

	user, err := h.store.FindByUsername(r.Context(), userName)
	if err != nil {
		log.Logger(r.Context()).Infof("Find error: %v", err)
		util.WriteJSON(w, http.StatusNotFound, []*Response{})
		return
	}

	util.WriteJSON(w, http.StatusOK, user.ToResponse())
}
