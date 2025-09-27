package handler

import "net/http"

type (
	ApiHandler interface {
		Register(mux *http.ServeMux)
	}
)
