package internal

import (
	"encoding/json"
	"net/http"
)

func RenderInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	renderJsonErrorPayload(w, r, NewInternalServerErrResponse(err), http.StatusInternalServerError)
}

func RenderBadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	renderJsonErrorPayload(w, r, NewBadRequestErrResponse(err), http.StatusBadRequest)
}

func RenderNotFoundError(w http.ResponseWriter, r *http.Request) {
	renderJsonErrorPayload(w, r, NotFoundErrResponse, http.StatusNotFound)
}

func RenderMethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func renderJsonErrorPayload(w http.ResponseWriter, _ *http.Request, value any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(value)
	if err != nil {
		panic(err)
	}
}
