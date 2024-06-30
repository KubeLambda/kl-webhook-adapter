package internal

import (
	"encoding/json"
	"net/http"
)

var AppVersion = VersionRest{
	Service: "rest-net/http",
	Version: "0.1.0",
	Build:   "1",
}

func VersionRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getVersion(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(AppVersion); err != nil {
		RenderInternalServerError(w, r, err)
	}
}
