package utils

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// IP gets ip from r.Header's X-FORWARDED-FOR or r.RemoteAddr
func IP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// IsMethod check if the method meets the requirement
// and response 405 Method Not Allowed if not matched
func IsMethod(m string, w http.ResponseWriter, r *http.Request) bool {
	logrus.Infoln("[utils.IsMethod]\t accept", IP(r), r.Method, r.URL)
	if r.Method != m {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
