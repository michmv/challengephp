package lib

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type ServerConfig[S, C, R any] struct {
	Storage            S
	RegistrationRoutes func(*RouterManager[S, C, R])
	MakeContext        func(S, http.ResponseWriter, *http.Request) C
	CloseContext       func(C)
	InitServerEngine   func(S) (*chi.Mux, Error)
}
