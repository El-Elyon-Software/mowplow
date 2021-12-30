package controllers

import (
	// "mowplow/api/v1/dal"
	// "mowplow/api/v1/models"

	"github.com/go-chi/chi/v5"
)

func EndCustomerServiceRoutes() *chi.Mux {
	router := chi.NewRouter()
	// ecs := models.NewEndCustomerService(dal.DB{})
	// router.Post("/", ecs.create)
	// router.Get("/{ID}", ecs.read)
	// router.Get("/", ecs.readFilter)
	// router.Put("/{ID}", ecs.update)
	// router.Delete("/{ID}", ecs.delete)
	return router
}
