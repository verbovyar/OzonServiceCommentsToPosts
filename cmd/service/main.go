package main

import (
	"log"
	"net/http"
	"ozonProject/graph"
	"ozonProject/internal/pubsub"
	"ozonProject/internal/service"
	databases "ozonProject/internal/storage/dataBases"
	"ozonProject/pkg/postgres"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {

	pool := postgres.New("postgres://postgres:Verbov323213@localhost:5432/OzonPostsCommentsDB")
	repo := databases.NewPostgresRepository(pool.Pool)
	svc := service.New(repo)
	bus := pubsub.New()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{Service: svc, Bus: bus},
	}))
	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("listening on %s", ":8080")
	log.Println("Playground:  http://localhost:8080/playground")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
