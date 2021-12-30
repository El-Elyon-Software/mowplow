// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type EndCustomer struct {
	ID           *int    `json:"id"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	BusinessName *string `json:"businessName"`
	Address1     string  `json:"address1"`
	Address2     *string `json:"address2"`
	PostalCode   string  `json:"postalCode"`
	Email        string  `json:"email"`
	Mobile       string  `json:"mobile"`
	DateAdded    *string `json:"dateAdded"`
	DateModified *string `json:"dateModified"`
}

type EndCustomerInput struct {
	ID           *int    `json:"id"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	BusinessName *string `json:"businessName"`
	Address1     string  `json:"address1"`
	Address2     *string `json:"address2"`
	PostalCode   string  `json:"postalCode"`
	Email        string  `json:"email"`
	Mobile       string  `json:"mobile"`
	DateAdded    *string `json:"dateAdded"`
	DateModified *string `json:"dateModified"`
}
