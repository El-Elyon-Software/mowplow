package service

import (
	"net/http"
	"strconv"
	"strings"
	"../dal"
	e "../errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/pkg/errors"
	"fmt"
)

//The choice has been made to make the business entities
//contain the CRUD operations and business logic governing
//CRUD. This can be modified later if need be.
func Routes() *chi.Mux {
	/* Putting the instantiation
	 of these objects in here for now. There is probably a better
	way to do this but I don't know what that is at this point.*/

	db := dal.DB{
		DBType:     "",
		DBName:     "",
		DBUser:     "",
		DBPassword: ""}
	db.NewDB()
	srv := Service{dal: db}

	router := chi.NewRouter()
	router.Post("/", srv.Create)
	router.Get("/{ID}", srv.Read)
	router.Get("/", srv.ReadFilter)
	router.Put("/{ID}", srv.Update)
	router.Delete("/{ID}", srv.Delete)
	router.Patch("/{ID}", srv.Update)
	return router
}

type Service struct {
	ID           int64  `json:"id"`
	ServiceName  string `json:"serviceName"`
	Description  string `json:"description"`
	DateAdded    string `json:"dateAdded"`
	DateModified string `json:"dateModified"`
	dal          dal.DB
}

type GeneralResponse struct {
	MSG string `json:"msg"`
	ID  int64  `json:"id"`
}

func (srv *Service) Create(rw http.ResponseWriter, r *http.Request) {
	if srv.bindData(rw, r) != nil {
		return
	}

	err := srv.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer srv.dal.CloseDB()

	stmt := `INSERT INTO service (
				service_name
				,description) 
			VALUES (?, ?);`
	
	res, err := srv.dal.DB.Exec(stmt, srv.ServiceName, srv.Description)

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
		srv.ID = id
	}
	render.JSON(rw, r, srv)
}

func (srv *Service) Read(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	err = srv.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer srv.dal.CloseDB()

	stmt := `SELECT 
				service_id
				,service_name
				,description
				,date_added
				,date_modified
			FROM 
				service 
			WHERE 
				service_id=?
				AND date_deleted IS NULL;`

	err = srv.dal.DB.QueryRow(stmt, id).Scan(
		&srv.ID, &srv.ServiceName, &srv.Description, 
		&srv.DateAdded, &srv.DateModified)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	render.JSON(rw, r, srv)
}

func (srv *Service) ReadFilter(rw http.ResponseWriter, r *http.Request) {
	wc, vals := dal.ParseQueryStringParams(r.URL.RawQuery)

	stmt := `SELECT 
				service_id
				,service_name
				,description
				,date_added
				,date_modified
			FROM 
				service 
			WHERE 
				` + strings.Join(wc, " ") + `;`
	fmt.Println(stmt)
	err := srv.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer srv.dal.CloseDB()

	rows, err := srv.dal.DB.Query(stmt, vals...)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer rows.Close()

	var ecs []*Service
	for rows.Next() {
		var s Service
		rows.Scan(&s.ID, &s.ServiceName, &s.Description, &s.DateAdded, &s.DateModified)
		ecs = append(ecs, &s)
	}

	render.JSON(rw, r, ecs)
}

func (srv *Service) Update(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	srv.ID = id

	if srv.bindData(rw, r) != nil {
		return
	}

	err = srv.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer srv.dal.CloseDB()

	stmt := `UPDATE 
				service 
			SET
				service_name=?
				,description=?
				,date_modified=NOW()
			WHERE
				service_id=?;`

	_, err = srv.dal.DB.Exec(stmt, srv.ServiceName, srv.Description, srv.ID)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	srv.Read(rw, r)
}

func (srv *Service) Delete(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	err = srv.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer srv.dal.CloseDB()

	stmt := `UPDATE 
				service 
			SET
				date_deleted=NOW()
			WHERE
				service_id=?;`

	_, err = srv.dal.DB.Exec(stmt, id)

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

func (srv *Service) bindData(rw http.ResponseWriter, r *http.Request) error {
	if err := render.Bind(r, srv); err != nil {
		e.HandleError(rw, r, err)
		return err
	}

	return nil
}

//ec implements the Bind interface, making this the Binder method called from render.binder
//On binding, request params are validated
func (srv *Service) Bind(r *http.Request) error {
	if srv.ServiceName == "" || len(srv.ServiceName) < 1 {
		return errors.New("serviceName is required and must be at least one characters.")
	}
	if srv.Description == "" || len(srv.Description) < 1 {
		return errors.New("description is required and must be at least one characters.")
	}
	
	return nil
}
