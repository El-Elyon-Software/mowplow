package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.RedirectSlashes,
		render.SetContentType(render.ContentTypeJSON),
	)

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/endCustomer", EndCustomerRoutes())
		// r.Mount("/service", buildEndPoint("/service").Routes())
		// r.Mount("/serviceProvider", buildEndPoint("/serviceProvider").Routes())
		// r.Mount("/serviceProviderService", buildEndPoint("/serviceProviderService").Routes())
		// r.Mount("/endCustomerService", buildEndPoint("/endCustomerService").Routes())
	})

	return r
}

type GeneralResponse struct {
	MSG string `json:"msg"`
	ID  int64  `json:"id"`
}
