package endPoint

import (
	"github.com/go-chi/chi"
)

type EndPoint interface {
	Routes() *chi.Mux
}