package models

import (
	"mowplow/api/v1/dal"
	e "mowplow/api/v1/errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

type ServiceProviderService struct {
	ID                int64  `json:"id"`
	ServiceProviderID int64  `json:"serviceProviderID"`
	ServiceID         int64  `json:"serviceID"`
	Description       string `json:"description"`
	DateAdded         string `json:"dateAdded"`
	DateModified      string `json:"dateModified"`
	dal               dal.DB
}

func (sp *ServiceProviderService) New() *ServiceProviderService {
	db := dal.DB{
		DBType:     "",
		DBName:     "",
		DBUser:     "",
		DBPassword: ""}
	//db.NewDB()
	sps := ServiceProviderService{dal: db}
	return &sps
}

func (sps *ServiceProviderService) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", sps.create)
	router.Get("/{ID}", sps.read)
	router.Get("/", sps.readFilter)
	router.Put("/{ID}", sps.update)
	router.Delete("/{ID}", sps.delete)
	return router
}

func (sps *ServiceProviderService) create(rw http.ResponseWriter, r *http.Request) {
	if sps.bindData(rw, r) != nil {
		return
	}

	err := sps.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer sps.dal.CloseDB()

	stmt := `INSERT INTO service_provider_service (
				service_provider_id
				,service_id
				,description) 
			VALUES (?, ?, ?);`

	res, err := sps.dal.DB.Exec(stmt, sps.ServiceProviderID, sps.ServiceID, sps.Description)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	if id > 0 {
		sps.ID = id
	}
	render.JSON(rw, r, sps)
}

func (sps *ServiceProviderService) read(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	err = sps.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer sps.dal.CloseDB()

	stmt := `SELECT 
				service_provider_service_id
				,service_provider_id
				,service_id
				,description
				,date_added
				,date_modified
			FROM 
				service_provider_service 
			WHERE 
				service_id=?
				AND date_deleted IS NULL;`

	err = sps.dal.DB.QueryRow(stmt, id).Scan(
		&sps.ID, &sps.ServiceProviderID, &sps.ServiceID,
		&sps.Description, &sps.DateAdded, &sps.DateModified)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	render.JSON(rw, r, sps)
}

func (sps *ServiceProviderService) readFilter(rw http.ResponseWriter, r *http.Request) {
	wc, vals := dal.ParseQueryStringParams(r.URL.RawQuery)

	stmt := `SELECT 
				service_provider_service_id
				,service_provider_id
				,service_id
				,description
				,date_added
				,date_modified
			FROM 
				service_provider_service 
			WHERE 
				` + strings.Join(wc, " ") + `;`

	err := sps.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer sps.dal.CloseDB()

	rows, err := sps.dal.DB.Query(stmt, vals...)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer rows.Close()

	var spss []*ServiceProviderService
	for rows.Next() {
		var s ServiceProviderService
		rows.Scan(&s.ID, &s.ServiceProviderID, &s.ServiceID, &s.Description, &s.DateAdded, &s.DateModified)
		spss = append(spss, &s)
	}

	render.JSON(rw, r, spss)
}

func (sps *ServiceProviderService) update(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	sps.ID = id

	if sps.bindData(rw, r) != nil {
		return
	}

	err = sps.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer sps.dal.CloseDB()

	stmt := `UPDATE 
				service_provider_service 
			SET
				service_provider_id=?
				,service_id=?
				,description=?
				,date_modified=NOW()
			WHERE
				service_provider_service_id=?;`

	_, err = sps.dal.DB.Exec(stmt, sps.ServiceProviderID, sps.ServiceID, sps.Description, sps.ID)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	sps.read(rw, r)
}

func (sps *ServiceProviderService) delete(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	err = sps.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer sps.dal.CloseDB()

	stmt := `UPDATE 
				service_provider_service 
			SET
				date_deleted=NOW()
			WHERE
				service_provider_service_id=?;`

	_, err = sps.dal.DB.Exec(stmt, id)

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

func (sps *ServiceProviderService) bindData(rw http.ResponseWriter, r *http.Request) error {
	if err := render.Bind(r, sps); err != nil {
		e.HandleError(rw, r, err)
		return err
	}

	return nil
}

func (sps *ServiceProviderService) Bind(r *http.Request) error {
	if sps.ServiceProviderID <= 0 {
		return errors.New("serviceProviderID is required and must be at least one characters.")
	}

	if sps.Description == "" || len(sps.Description) < 1 {
		return errors.New("description is required and must be at least one characters.")
	}

	return nil
}
