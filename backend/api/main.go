package api

import (
	"net/http"

	"github.com/fumiama/paper-manager/backend/utils"
)

// Handler serves all backend /api call
func Handler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsMethod("GET", w, r) {
		return
	}
	http.Error(w, "404 Not Found", http.StatusNotFound)
}
