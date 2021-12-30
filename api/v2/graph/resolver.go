package graph

import (
	"mowplow/api/v2/dal"
	// "github.com/99designs/gqlgen/example/federation/reviews/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *dal.DB
}
