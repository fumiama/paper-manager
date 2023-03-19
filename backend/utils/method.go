package utils

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// IsMethod check if the method meets the requirement
// and response 405 Method Not Allowed if not matched
func IsMethod(m string, w http.ResponseWriter, r *http.Request) bool {
	logrus.Infoln("[utils.IsMethod] accept", r.RemoteAddr, r.Method, r.URL)
	if r.Method != m {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
