package services

import (
	"mowplow/api/v1/dal"
	"mowplow/api/v1/models"
	"strings"
)

func RetrieveEndCustomersWithFilters(wc []string, vals []interface{}, db *dal.DB) ([]models.EndCustomer, error) {
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

	err := db.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.CloseDB()

	rows, err := db.DB.Query(stmt, vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ecs []models.EndCustomer
	for rows.Next() {
		var e models.EndCustomer
		rows.Scan(
			&e.ID, &e.FirstName, &e.LastName, &e.BusinessName,
			&e.Address1, &e.Address2, &e.PostalCode,
			&e.Email, &e.Mobile, &e.DateAdded, &e.DateModified)
		ecs = append(ecs, e)
	}

	return ecs, nil
}
