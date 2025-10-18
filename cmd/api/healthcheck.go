package main

import (
	"dijam-ecommerce/shared"
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Status string `json:"status"`
	}{
		Method: r.Method,
		Path:   r.URL.RequestURI(),
		Status: "available",
	}
	shared.WriteToJSON(w, http.StatusOK, shared.Envelope{"health": data}, nil)
}
