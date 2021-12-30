package models

import (
	"mowplow/api/v1/dal"
)

type endCustomerService struct {
	ID                       int64   `json:"id"`
	EndCustomerID            int64   `json:"endCustomerID"`
	ServiceProviderServiceID int64   `json:"serviceProviderServiceID"`
	Description              string  `json:"description"`
	EstimatedJobLength       float64 `json:"estimatedJobLength"`
	ContractStartDate        string  `json:"contractStartDate"`
	ContractEndDate          string  `json:"contractEndDate"`
	DateAdded                string  `json:"dateAdded"`
	DateModified             string  `json:"dateModified"`
	dal                      dal.DB
}

func NewEndCustomerService(db dal.DB) *endCustomerService {
	ecs := endCustomerService{dal: db}
	return &ecs
}
