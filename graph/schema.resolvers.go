package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/bellwood4486/gqlgen-todos/db"
	"github.com/bellwood4486/gqlgen-todos/graph/generated"
	"github.com/bellwood4486/gqlgen-todos/graph/model"
)

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	res := db.LogAndQuery(r.Conn, "SELECT id, text, user_id FROM todos")
	defer res.Close()

	var todos []*model.Todo
	for res.Next() {
		var todo model.Todo
		if err := res.Scan(&todo.ID, &todo.Text, &todo.UserID); err != nil {
			panic(err)
		}
		todos = append(todos, &todo)
	}

	return todos, nil
}

func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	res := db.LogAndQuery(r.Conn, "SELECT id, name FROM users WHERE id = $1", obj.UserID)
	defer res.Close()

	if !res.Next() {
		return nil, nil
	}
	var user model.User
	if err := res.Scan(&user.ID, &user.Name); err != nil {
		panic(err)
	}
	return &user, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Todo returns generated.TodoResolver implementation.
func (r *Resolver) Todo() generated.TodoResolver { return &todoResolver{r} }

type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
