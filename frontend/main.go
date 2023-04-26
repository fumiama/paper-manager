package frontend

import (
	"net/http"

	"github.com/fumiama/paper-manager/frontend/vben"
)

// StaticHandler serves contents in frontend
var StaticHandler = http.FileServer(vben.Distribution)
