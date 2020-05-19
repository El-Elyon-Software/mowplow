package main

import (
	"log"
	"net/http"
	"./endCustomer"
	"./service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/render"
	"time"
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
		r.Mount("/endCustomer", endCustomer.NewEndCustomer())
		r.Mount("/service", service.NewService())
	})

	return r
}

func main() {
	time.Sleep(20)
	router := Routes()
	log.Fatal(http.ListenAndServe(":8080", router))
}