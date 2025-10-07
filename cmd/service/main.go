package main

import (
	"log"
	"net/http"
	"ozonProject/graph"
	"ozonProject/internal/service"
	databases "ozonProject/internal/storage/dataBases"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	st := databases.NewInMemRepository()
	svc := service.New(st)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{Service: svc},
	}))

	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("listening on %s", ":8080")
	log.Println("Playground:  http://localhost:8080/playground")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
