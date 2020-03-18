//go:generate go run github.com/vektah/dataloaden UserLoader int32 *github.com/bellwood4486/gqlgen-todos/graph/model.User

package dataloader

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bellwood4486/gqlgen-todos/db"

	"github.com/bellwood4486/gqlgen-todos/graph/model"
)

const loadersKey = "dataloaders"

type Loaders struct {
	UserById UserLoader
}

func Middleware(conn *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			UserById: UserLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int32) ([]*model.User, []error) {
					placeholders := make([]string, len(ids))
					args := make([]interface{}, len(ids))
					for i := 0; i < len(ids); i++ {
						placeholders[i] = fmt.Sprintf("$%d", i+1)
						args[i] = ids[i]
					}

					res := db.LogAndQuery(conn,
						"SELECT id, name FROM users WHERE id IN ("+strings.Join(placeholders, ",")+")",
						args...)
					defer res.Close()

					usersById := map[int32]*model.User{}
					for res.Next() {
						user := model.User{}
						err := res.Scan(&user.ID, &user.Name)
						if err != nil {
							panic(err)
						}
						usersById[user.ID] = &user
					}

					users := make([]*model.User, len(ids))
					for i, id := range ids {
						users[i] = usersById[id]
						i++
					}

					return users, nil
				},
			}})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
