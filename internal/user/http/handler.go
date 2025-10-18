package userHTTPHandler

import (
	users "dijam-ecommerce/internal/user"
	"dijam-ecommerce/pkg/apierrors"
	"dijam-ecommerce/shared"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service   *users.UserService
	validator *validator.Validate
}

func NewHandler(s *users.UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{
		service:   s,
		validator: v,
	}
}

func formatValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)

	// Type assert to get the validator.ValidationErrors
	if vErrs, ok := err.(validator.ValidationErrors); ok {
		for _, vErr := range vErrs {
			// Use the JSON field name (or just Field() if you don't use JSON tags)
			field := strings.ToLower(vErr.Field())

			// Provide user-friendly messages
			switch vErr.Tag() {
			case "required":
				errorsMap[field] = "this field is required"
			case "email":
				errorsMap[field] = "must be a valid email address"
			case "min":
				errorsMap[field] = "this field is too short"
			case "max":
				errorsMap[field] = "this field is too long"
			default:
				errorsMap[field] = "this field is invalid"
			}
		}
	}
	return errorsMap
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name" validate:"required,min=2"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=72"`
	}
	err := shared.ReadJSON(w, r, &req)
	if err != nil {
		apierrors.BadRequest(w, r, err, "invalid JSON Request")
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		validationErrors := formatValidationErrors(err)
		apierrors.ValidationFailed(w, r, validationErrors)
		return
	}

	err = h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, users.ErrDuplicateEmail) {
			apierrors.Conflict(w, r, err, "email already exists")
			return
		}
		apierrors.ServerError(w, r, err)
		return
	}
	_ = shared.WriteToJSON(w, http.StatusCreated, shared.Envelope{"message": "user successfully created"}, nil)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	err := shared.ReadJSON(w, r, &req)
	if err != nil {
		apierrors.BadRequest(w, r, err, "invalid JSON payload")
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		validationErrors := formatValidationErrors(err)
		apierrors.ValidationFailed(w, r, validationErrors)
		return
	}

	err = h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, users.ErrInvalidCredentials) {
			apierrors.Unauthorized(w, r, err)
			return
		}
		apierrors.ServerError(w, r, err)
		return
	}

	_ = shared.WriteToJSON(w, http.StatusCreated, shared.Envelope{"message": "user successfully logged in"}, nil)
}
