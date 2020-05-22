package main

import (
	"log"
	"net/http"
	"./endCustomer"
	ep "./endPoint"
	"./service"
	"./serviceProvider"
	"./serviceProviderService"
	"./endCustomerService"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/render"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		render.SetContentType(render.ContentTypeJSON),
	)

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/endCustomer", buildEndPoint("/endCustomer").Routes())
		r.Mount("/service", buildEndPoint("/service").Routes())
		r.Mount("/serviceProvider", buildEndPoint("/serviceProvider").Routes())
		r.Mount("/serviceProviderService", buildEndPoint("/serviceProviderService").Routes())
		r.Mount("/endCustomerService", buildEndPoint("/endCustomerService").Routes())
	})

	return r
}

func buildEndPoint(ep string) ep.EndPoint {
	switch s := ep; s {
	case "/endCustomer":
		return endCustomer.New()
	case "/service":
		return service.New()
	case "/serviceProvider":
		return serviceProvider.New()
	case "/serviceProviderService":
		return serviceProviderService.New()
	case "/endCustomerService":
		return endCustomerService.New()
	}
	return nil
}

func main() {
	router := Routes()
	log.Fatal(http.ListenAndServe(":8080", router))
}