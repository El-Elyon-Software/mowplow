package main

import (
	"log"
	"net/http"
	"./endCustomer"
	ep "./endPoint"
	"./service"
	"./serviceProvider"
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
	})

	return r
}

func buildEndPoint(ep string) ep.EndPoint {
	switch s := ep; s {
	case "/endCustomer":
		return endCustomer.NewEndCustomer()
	case "/service":
		return service.NewService()
	case "/serviceProvider":
		return serviceProvider.NewServiceProvider()
	}
	return nil
}

func main() {
	// time.Sleep(20)
	router := Routes()
	log.Fatal(http.ListenAndServe(":8080", router))
}