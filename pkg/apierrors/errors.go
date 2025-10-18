package apierrors

import (
	"dijam-ecommerce/shared"
	"log/slog"
	"net/http"
)

func logError(r *http.Request, err error, message string, details ...any) {
	args := []any{
		"method", r.Method, "path", r.URL.Path, "error", err.Error(),
	}

	args = append(args, details...)
	slog.Error(message, args)
}


func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err, "internal server error")
	msg := shared.Envelope{"error": "the server encountered a problem and could not process your request"}
	_ = shared.WriteToJSON(w, http.StatusInternalServerError, msg, nil)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error, message string) {
	logError(r, err, "bad request", "details", message)
	msg := shared.Envelope{"error": message}
	_ = shared.WriteToJSON(w, http.StatusBadRequest, msg, nil)
}

func ValidationFailed(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	slog.Warn("validation failed", "method", r.Method, "path", r.URL.Path, "errors", errors)
	msg := shared.Envelope{"errors": errors}
	_ = shared.WriteToJSON(w, http.StatusUnprocessableEntity, msg, nil)
}

func Conflict(w http.ResponseWriter, r *http.Request, err error, message string) {
	logError(r, err, "conflict", "details", message)
	msg := shared.Envelope{"error": message}
	_ = shared.WriteToJSON(w, http.StatusConflict, msg, nil)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("unauthorized attempt", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	msg := shared.Envelope{"error": "invalid authentication credentials"}
	_ = shared.WriteToJSON(w, http.StatusUnauthorized, msg, nil)
}
