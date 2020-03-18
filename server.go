package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bellwood4486/gqlgen-todos/dataloader"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bellwood4486/gqlgen-todos/graph"
	"github.com/bellwood4486/gqlgen-todos/graph/generated"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func main() {
	// docker run --rm --name dataloader-example -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_DB=dataloader_example -d postgres:9.6
	db, err := sql.Open("postgres", "user=root dbname=dataloader_example sslmode=disable")
	if err != nil {
		panic(err)
	}

	mustExec(db, "DROP TABLE IF EXISTS users")
	mustExec(db, "CREATE TABLE users (id serial primary key, name varchar(255))")
	mustExec(db, "DROP TABLE IF EXISTS todos")
	mustExec(db, "CREATE TABLE todos (id serial primary key, text varchar(255), user_id int)")
	for i := 1; i <= 5; i++ {
		mustExec(db, "INSERT INTO users (name) VALUES ($1)", fmt.Sprintf("user %d", i))
	}
	for i := 1; i <= 20; i++ {
		mustExec(db, "INSERT INTO todos (text, user_id) VALUES ($1, $2)", fmt.Sprintf("Todo %d", i), (i%5)+1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Conn: db,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", dataloader.Middleware(db, srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func mustExec(db *sql.DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		panic(err)
	}
}
