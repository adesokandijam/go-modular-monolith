package userHTTPHandler

import (
	"dijam-ecommerce/internal/user"
	"dijam-ecommerce/shared"
	"errors"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	service *user.UserService
}

func NewHandler(s *user.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := shared.ReadJSON(w, r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, user.ErrDuplicateEmail) {
			shared.WriteToJSON(w, http.StatusConflict, shared.Envelope{"error": "email already exists"}, nil)
			slog.Error("email already exist")
			return
		}
		shared.WriteToJSON(w, http.StatusInternalServerError, shared.Envelope{"error": "problem with creating user"}, nil)
		slog.Error("problem with creating the user", "err", err.Error())
		return
	}
	_ = shared.WriteToJSON(w, http.StatusCreated, shared.Envelope{"message": "user successfully created"}, nil)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := shared.ReadJSON(w, r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		shared.WriteToJSON(w, http.StatusInternalServerError, shared.Envelope{"error": "incorrect credentials"}, nil)
		slog.Error("unable to login user", "err", err.Error())
		return
	}

	_ = shared.WriteToJSON(w, http.StatusCreated, shared.Envelope{"message": "user successfully logged in"}, nil)
}
