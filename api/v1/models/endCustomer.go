package models

import (
	"mowplow/api/v1/dal"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
)

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
	dal          *dal.DB
}

func NewEndCustomer(db *dal.DB) *EndCustomer {
	ec := EndCustomer{dal: db}
	return &ec
}

/*
Used for creating new end customers customers.
Returns an error
*/
func (ec *EndCustomer) Create() error {
	err := ec.dal.OpenDB()
	if err != nil {
		return err
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
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	if id > 0 {
		ec.ID = id
	}

	return nil
}

/*
Retrieves the end customer object associated with the given ID
*/
func (ec *EndCustomer) Retrieve(id int64) error {
	err := ec.dal.OpenDB()
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (ec *EndCustomer) Update() error {
	err := ec.dal.OpenDB()
	if err != nil {
		return err
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
		ec.Address1, ec.Address2, ec.PostalCode, ec.Email, ec.Mobile, ec.ID)
	if err != nil {
		return err
	}

	return ec.Retrieve(ec.ID)
}

func (ec *EndCustomer) Delete(id int64) error {
	err := ec.dal.OpenDB()
	if err != nil {
		return err
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
		return err
	}
	
	return nil
}

//EndCustomer implements the Bind interface, making this the Binder method called from render.Binder
//On binding, EndCustomer fields are validated
func (ec EndCustomer) Bind(r *http.Request) error {
	if ec.FirstName == "" || len(ec.FirstName) == 0 {
		return errors.New("firstName is required and must be at least one characters.")
	}
	if ec.LastName == "" || len(ec.LastName) == 0 {
		return errors.New("lastName is required and must be at least one characters.")
	}
	if ec.Address1 == "" || len(ec.Address1) == 0 {
		return errors.New("address1 is required and must be at least four characters.")
	}
	if ec.PostalCode == "" || len(ec.PostalCode) < 5 {
		return errors.New("postalCode is required and must be at least five characters.")
	}
	if len(ec.PostalCode) > 10 {
		return errors.New("postalCode is too long, must be less than 11 characters.")
	}

	if ec.Email == "" || len(ec.Email) == 0 {
		return errors.New("email is required.")
	}
	re := regexp.MustCompile(`(?mi)[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,6}`)
	match := re.FindAllString(ec.Email, -1)
	if match == nil {
		return errors.New("Invalid email.")
	}

	if ec.Mobile == "" || len(ec.Mobile) == 0 {
		return errors.New("mobile is required.")
	}

	// This is only for North American numbers
	re = regexp.MustCompile(`^\D?(\d{3})\D?\D?(\d{3})\D?(\d{4})$`)
	match = re.FindAllString(ec.Mobile, -1)
	if match == nil {
		return errors.New("Invalid mobile.")
	}

	return nil
}
