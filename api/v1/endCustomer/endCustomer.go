package endCustomer

import (
	"net/http"
	"regexp"
	"strconv"

	"../dal"
	e "../errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/pkg/errors"
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
	ec := EndCustomer{dal: db}

	router := chi.NewRouter()
	router.Post("/", ec.Create)
	router.Get("/{ID}", ec.Read)
	router.Get("/", ec.ReadFilter)
	router.Put("/{ID}", ec.Update)
	router.Delete("/{ID}", ec.Delete)
	return router
}

type EndCustomer struct {
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

func (ec *EndCustomer) Create(rw http.ResponseWriter, r *http.Request) {
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

func (ec *EndCustomer) Read(rw http.ResponseWriter, r *http.Request) {
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

func (ec *EndCustomer) ReadFilter(rw http.ResponseWriter, r *http.Request) {
	wc, vals := dal.ParseQueryStringParams(r.URL.Query())

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
				` + wc[0] + `;`

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

	var ecs []*EndCustomer
	for rows.Next() {
		var singel_ec EndCustomer
		rows.Scan(
			&singel_ec.ID, &singel_ec.FirstName, &singel_ec.LastName,
			&singel_ec.BusinessName,
			&singel_ec.Address1, &singel_ec.Address2, &singel_ec.PostalCode,
			&singel_ec.Email, &singel_ec.Mobile,
			&singel_ec.DateAdded, &singel_ec.DateModified)

		ecs = append(ecs, &singel_ec)
	}

	render.JSON(rw, r, ecs)
}

func (ec *EndCustomer) Update(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	ec.ID = id

	if ec.bindData(rw, r) != nil {
		return
	}

	err = ec.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, err)
		return
	}
	defer ec.dal.CloseDB()

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
				end_customer_id=?;`

	_, err = ec.dal.DB.Exec(stmt, ec.FirstName, ec.LastName, ec.BusinessName,
		ec.Address1, ec.Address2, ec.PostalCode, ec.Email, ec.Mobile,
		ec.ID)

	if err != nil {
		e.HandleError(rw, r, err)
		return
	}

	ec.Read(rw, r)
}

func (ec *EndCustomer) Delete(rw http.ResponseWriter, r *http.Request) {
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

	stmt := `UPDATE 
				end_customer 
			SET
				date_deleted=NOW()
			WHERE
				end_customer_id=?;`

	_, err = ec.dal.DB.Exec(stmt, id)

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

func (ec *EndCustomer) DeleteFilter(rw http.ResponseWriter, r *http.Request) {

}

func (ec *EndCustomer) bindData(rw http.ResponseWriter, r *http.Request) error {
	if err := render.Bind(r, ec); err != nil {
		e.HandleError(rw, r, err)
		return err
	}

	return nil
}

func (ec *EndCustomer) Bind(r *http.Request) error {
	if ec.FirstName == "" || len(ec.FirstName) < 1 {
		return errors.New("firstName is required and must be at least one characters.")
	}
	if ec.LastName == "" || len(ec.LastName) < 1 {
		return errors.New("lastName is required and must be at least one characters.")
	}
	if ec.Address1 == "" || len(ec.Address1) < 4 {
		return errors.New("address1 is required and must be at least four characters.")
	}
	if ec.PostalCode == "" || len(ec.PostalCode) < 4 {
		return errors.New("postalCode is required and must be at least four characters.")
	}

	re := regexp.MustCompile(`(?mi)[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,6}`)
	match := re.FindAllString(ec.Email, -1)
	if match == nil {
		return errors.New("Invalid email.")
	}

	// This is only for North American numbers
	re = regexp.MustCompile(`^\D?(\d{3})\D?\D?(\d{3})\D?(\d{4})$`)
	match = re.FindAllString(ec.Mobile, -1)
	if match == nil {
		return errors.New("Invalid mobile.")
	}

	return nil
}
