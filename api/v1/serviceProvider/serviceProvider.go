package serviceProvider

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"../dal"
	e "../errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/pkg/errors"
)

type ServiceProvider struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	BusinessName string `json:"businessName"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	PostalCode   string `json:"postalCode"`
	Email        string `json:"email"`
	Mobile       string `json:"mobile"`
	DateAdded    string `json:"dateAdded"`
	DateModified string `json:"dateModified"`
	dal          dal.DB
}

type GeneralResponse struct {
	MSG string `json:"msg"`
	ID  int64  `json:"id"`
}

func NewServiceProvider() *ServiceProvider {
	db := dal.DB{
		DBType:     "",
		DBName:     "",
		DBUser:     "",
		DBPassword: ""}
	db.NewDB()
	ec := ServiceProvider{dal: db}
	return &ec
}

func (sp *ServiceProvider) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", sp.create)
	router.Get("/{ID}", sp.read)
	router.Get("/", sp.readFilter)
	router.Put("/{ID}", sp.update)
	router.Delete("/{ID}", sp.delete)
	router.Patch("/{ID}", sp.update)
	return router
}

func (ec *ServiceProvider) create(rw http.ResponseWriter, r *http.Request) {
	if ec.bindData(rw, r) != nil {
		return
	}

	err := ec.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer ec.dal.CloseDB()

	stmt := `INSERT INTO end_customer (
				first_name
				,last_name
				, business_name
				, address_1
				, address_2
				, postal_code
				, email
				, mobile) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	res, err := ec.dal.DB.Exec(stmt, ec.FirstName, ec.LastName, ec.BusinessName,
		ec.Address1, ec.Address2, ec.PostalCode, ec.Email, ec.Mobile)

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
		ec.ID = id
	}
	render.JSON(rw, r, ec)

}

func (ec *ServiceProvider) read(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	err = ec.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer ec.dal.CloseDB()

	stmt := `SELECT 
				end_customer_id
				,first_name
				,last_name
				,business_name
				,address_1
				,address_2
				,postal_code
				,email
				,mobile
				,date_added
				,date_modified
			FROM 
				end_customer 
			WHERE 
				end_customer_id=?
				AND date_deleted IS NULL;`

	err = ec.dal.DB.QueryRow(stmt, id).Scan(
		&ec.ID, &ec.FirstName, &ec.LastName, &ec.BusinessName,
		&ec.Address1, &ec.Address2, &ec.PostalCode, &ec.Email, &ec.Mobile,
		&ec.DateAdded, &ec.DateModified)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	render.JSON(rw, r, ec)
}

func (ec *ServiceProvider) readFilter(rw http.ResponseWriter, r *http.Request) {
	wc, vals := dal.ParseQueryStringParams(r.URL.RawQuery)

	stmt := `SELECT 
				end_customer_id
				,first_name
				,last_name
				,business_name
				,address_1
				,address_2
				,postal_code
				,email
				,mobile
				,date_added
				,date_modified
			FROM 
				end_customer 
			WHERE 
				` + strings.Join(wc, " ") + `;`

	err := ec.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer ec.dal.CloseDB()

	rows, err := ec.dal.DB.Query(stmt, vals...)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer rows.Close()

	var sps []*ServiceProvider
	for rows.Next() {
		var sp ServiceProvider
		rows.Scan(
			&sp.ID, &sp.FirstName, &sp.LastName, &sp.BusinessName,
			&sp.Address1, &sp.Address2, &sp.PostalCode,
			&sp.Email, &sp.Mobile, &sp.DateAdded, &sp.DateModified)
		sps = append(sps, &sp)
	}

	render.JSON(rw, r, sps)
}

func (sp *ServiceProvider) update(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	sp.ID = id

	if sp.bindData(rw, r) != nil {
		return
	}

	err = sp.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer sp.dal.CloseDB()

	stmt := `UPDATE 
				end_customer 
			SET
				first_name=?
				,last_name=?
				,business_name=?
				,address_1=?
				,address_2=?
				,postal_code=?
				,email=?
				,mobile=?
				,date_modified=NOW()
			WHERE
				service_provider_id=?;`

	_, err = sp.dal.DB.Exec(stmt, sp.FirstName, sp.LastName, sp.BusinessName,
		sp.Address1, sp.Address2, sp.PostalCode, sp.Email, sp.Mobile, sp.ID)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	sp.read(rw, r)
}

func (sp *ServiceProvider) delete(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	err = sp.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer sp.dal.CloseDB()

	stmt := `UPDATE 
				service_provider 
			SET
				date_deleted=NOW()
			WHERE
				service_provider_id=?;`

	_, err = sp.dal.DB.Exec(stmt, id)

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

func (sp *ServiceProvider) bindData(rw http.ResponseWriter, r *http.Request) error {
	if err := render.Bind(r, sp); err != nil {
		e.HandleError(rw, r, err)
		return err
	}

	return nil
}

//ec implements the Bind interface, making this the Binder method called from render.binder
//On binding, request params are validated
func (sp *ServiceProvider) Bind(r *http.Request) error {
	if sp.FirstName == "" || len(sp.FirstName) < 1 {
		return errors.New("firstName is required and must be at least one characters.")
	}
	if sp.LastName == "" || len(sp.LastName) < 1 {
		return errors.New("lastName is required and must be at least one characters.")
	}
	if sp.Address1 == "" || len(sp.Address1) < 4 {
		return errors.New("address1 is required and must be at least four characters.")
	}
	if sp.PostalCode == "" || len(sp.PostalCode) < 4 {
		return errors.New("postalCode is required and must be at least four characters.")
	}

	re := regexp.MustCompile(`(?mi)[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,6}`)
	match := re.FindAllString(sp.Email, -1)
	if match == nil {
		return errors.New("Invalid email.")
	}

	// This is only for North American numbers
	re = regexp.MustCompile(`^\D?(\d{3})\D?\D?(\d{3})\D?(\d{4})$`)
	match = re.FindAllString(sp.Mobile, -1)
	if match == nil {
		return errors.New("Invalid mobile.")
	}

	return nil
}