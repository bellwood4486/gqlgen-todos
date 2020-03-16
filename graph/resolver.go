package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"database/sql"

	"github.com/bellwood4486/gqlgen-todos/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	todos []*model.Todo
	Conn  *sql.DB
}
