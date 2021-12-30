package controllers

import (
	"mowplow/api/v1/dal"
	e "mowplow/api/v1/errors"
	"mowplow/api/v1/models"
	"mowplow/api/v1/services"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type endCustomerController struct {
}

func EndCustomerRoutes() *chi.Mux {
	router := chi.NewRouter()
	ecc := endCustomerController{}

	router.Post("/", ecc.create)
	router.Get("/{ID}", ecc.retrieve)
	router.Get("/", ecc.retrieveWithFilter)
	router.Put("/", ecc.update)
	router.Delete("/{ID}", ecc.delete)
	return router
}

func (ecc endCustomerController) create(rw http.ResponseWriter, r *http.Request) {
	ec := models.NewEndCustomer(dal.NewDB())

	if err := render.Bind(r, ec); err != nil {
		e.HandleError(rw, r, err)
		return
	}

	if err := ec.Create(); err != nil {
		e.HandleError(rw, r, err)
		return
	}

	render.JSON(rw, r, ec)
}

func (ecc endCustomerController) retrieve(rw http.ResponseWriter, r *http.Request) {
	ec := models.NewEndCustomer(dal.NewDB())

	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	err = ec.Retrieve(id)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	render.JSON(rw, r, ec)
}

func (ecc endCustomerController) retrieveWithFilter(rw http.ResponseWriter, r *http.Request) {
	wc, vals := dal.ParseQueryStringParams(r.URL.RawQuery)

	ecs, err := services.RetrieveEndCustomersWithFilters(wc, vals, dal.NewDB())
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	render.JSON(rw, r, ecs)
}

func (ecc endCustomerController) update(rw http.ResponseWriter, r *http.Request) {
	ec := models.NewEndCustomer(dal.NewDB())

	if err := render.Bind(r, ec); err != nil {
		e.HandleError(rw, r, err)
		return
	}

	if err := ec.Update(); err != nil {
		e.HandleError(rw, r, err)
		return
	}

	render.JSON(rw, r, ec)
}

func (ecc endCustomerController) delete(rw http.ResponseWriter, r *http.Request) {
	ec := models.NewEndCustomer(dal.NewDB())

	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	err = ec.Delete(id)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	gr := &GeneralResponse{
		MSG: "The end customer was deleted",
		ID:  id,
	}

	render.JSON(rw, r, gr)
}
