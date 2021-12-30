package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"mowplow/api/v2/graph/generated"
	"mowplow/api/v2/graph/model"
)

func (r *mutationResolver) CreateEndCustomer(ctx context.Context, input model.EndCustomerInput) (*model.EndCustomer, error) {
	ec := &model.EndCustomer{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		BusinessName: input.BusinessName,
		Address1:     input.Address1,
		Address2:     input.Address2,
		PostalCode:   input.PostalCode,
		Email:        input.Email,
		Mobile:       input.Mobile,
	}

	err := r.DB.OpenDB()
	if err != nil {
		return nil, err
	}
	defer r.DB.CloseDB()

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

	res, err := r.DB.DB.Exec(stmt, ec.FirstName, ec.LastName, ec.BusinessName,
		ec.Address1, ec.Address2, ec.PostalCode, ec.Email, ec.Mobile)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	if id > 0 {
		id := int(id)
		ec.ID = &id
	}

	return ec, nil
}

func (r *mutationResolver) UpdateEndCustomer(ctx context.Context, input model.EndCustomerInput) (*model.EndCustomer, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteEndCustomer(ctx context.Context, input string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetEndCustomer(ctx context.Context, id int) (*model.EndCustomer, error) {
	err := r.DB.OpenDB()
	if err != nil {
		return nil, err
	}
	defer r.DB.CloseDB()

	ec := &model.EndCustomer{}

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

	err = r.DB.DB.QueryRow(stmt, id).Scan(
		&ec.ID, &ec.FirstName, &ec.LastName, &ec.BusinessName,
		&ec.Address1, &ec.Address2, &ec.PostalCode, &ec.Email, &ec.Mobile,
		&ec.DateAdded, &ec.DateModified)
	if err != nil {
		return nil, err
	}

	return ec, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
