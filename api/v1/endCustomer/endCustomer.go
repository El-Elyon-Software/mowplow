package endCustomer

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/pkg/errors"

	"../dal"
	e "../errors"
)

func Routes() *chi.Mux {
	// Hard coding for now. eventually will be replaced with dynamic customer specific creds.
	// Putting the instantiation of these objects in here for now.
	// There is probably a better way to do this but I don't know what that is at this point.
	dal := dal.DB{
		DBType:     "mysql",
		DBName:     "mowplow",
		DBUser:     "root",
		DBPassword: "saunya18!!",
	}
	ec := EndCustomer{dal: dal}

	router := chi.NewRouter()
	router.Post("/", ec.Create)
	router.Get("/{endCustomerID}", ec.Retrieve)
	router.Put("/", ec.Update)
	router.Delete("/", ec.Delete)
	return router
}

type EndCustomer struct {
	EndCustomerID int64  `json:"id"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	BusinessName  string `json:"businessName"`
	Address1      string `json:"address1"`
	Address2      string `json:"address2"`
	PostalCode    string `json:"postalCode"`
	Email         string `json:"email"`
	Mobile        string `json:"mobile"`
	DateAdded     string `json:"dateAdded"`
	DateModified  string `json:"dateModified"`
	dal           dal.DB
}

func (ec *EndCustomer) Create(rw http.ResponseWriter, r *http.Request) {
	ec.bindData(rw, r)

	err := ec.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, 500, "Internal Server Error", err)
		return
	}

	defer ec.dal.DB.Close()

	stmt := `INSERT INTO end_customer 
				(first_name,last_name, business_name, 
				address_1, address_2, postal_code, 
				email, mobile) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := ec.dal.DB.Exec(stmt, ec.FirstName, ec.LastName, ec.BusinessName,
		ec.Address1, ec.Address2, ec.PostalCode, ec.Email, ec.Mobile)

	if err != nil {
		e.HandleError(rw, r, 422, "Unprocessable Entity", err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		e.HandleError(rw, r, 500, "Internal Server Error", err)
		return
	}

	if id > 0 {
		ec.EndCustomerID = id
	}
	render.JSON(rw, r, ec)

}

func (ec *EndCustomer) Retrieve(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "endCustomerID"), 10, 64)
	if err != nil {
		e.HandleError(rw, r, 400, "Bad Request", err)
		return
	}

	err = ec.dal.OpenDB()
	if err != nil {
		e.HandleError(rw, r, 500, "Internal Server Error", err)
		return
	}
	defer ec.dal.DB.Close()

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
				end_customer_id=?`
	err = ec.dal.DB.QueryRow(stmt, id).Scan(
		&ec.EndCustomerID, &ec.FirstName, &ec.LastName, &ec.BusinessName, 
		&ec.Address1, &ec.Address2, &ec.PostalCode, &ec.Email, &ec.Mobile, 
		&ec.DateAdded, &ec.DateModified)

	if err != nil {
		e.HandleError(rw, r, 422, "Unprocessable Entity", err)
		return
	}

	render.JSON(rw, r, ec)
}

func (ec *EndCustomer) RetrieveFilter(rw http.ResponseWriter, r *http.Request) {
	ecs := []*EndCustomer{
		{
			EndCustomerID: 1,
			FirstName:     "BDP",
			LastName:      "In Da Place To Be",
			BusinessName:  "ROC City Coders",
			Address1:      "123 Roc City Blvd",
			Address2:      "Not your mom's place",
			PostalCode:    "12345-67",
			Email:         "bdp@indaplacetobe.com",
			Mobile:        "5855551212",
			DateAdded:     "01-01-2001",
			DateModified:  "12-01-2019",
		},
		{
			EndCustomerID: 2,
			FirstName:     "SDP",
			LastName:      "In Da Place To Be",
			BusinessName:  "ROC City Coders",
			Address1:      "321 Roc City Blvd",
			Address2:      "Not your dad's place",
			PostalCode:    "12345-67",
			Email:         "sdp@indaplacetobe.com",
			Mobile:        "5855551212",
			DateAdded:     "01-01-2001",
			DateModified:  "12-01-2019",
		},
	}
	render.JSON(rw, r, ecs)
}

func (ec *EndCustomer) Update(rw http.ResponseWriter, r *http.Request) {
	ec.bindData(rw, r)
}

func (ec *EndCustomer) Delete(rw http.ResponseWriter, r *http.Request) {
	ec.bindData(rw, r)
}

func (ec *EndCustomer) DeleteFilter(rw http.ResponseWriter, r *http.Request) {
	ec.bindData(rw, r)

}

func (ec *EndCustomer) bindData(rw http.ResponseWriter, r *http.Request) {
	if err := render.Bind(r, ec); err != nil {
		e.HandleError(rw, r, 422, "Unprocessable Entity", err)
	}
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
	if ec.Email == "" { //add regex check
		return errors.New("email is required.")
	}
	if ec.Mobile == "" { //add regex check
		return errors.New("mobile is required.")
	}

	return nil
}
