package main

import (
	"log"
	"net/http"
	"ozonProject/config"
	"ozonProject/graph"
	"ozonProject/internal/pubsub"
	"ozonProject/internal/service"
	databases "ozonProject/internal/storage/dataBases"
	"ozonProject/internal/storage/interfaces"
	"ozonProject/pkg/postgres"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	conf, err := config.LoadConfig("./config")
	if err != nil {
		println(err.Error())
	}

	var repo interfaces.StoreIface
	if conf.RepositoryType == "Postgres" {
		repo = startPostgres(conf.ConnectingString)
	} else if conf.RepositoryType == "InMemory" {
		repo = startInMemory()
	}

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

	log.Printf("listening on %s", conf.Port)
	log.Println("Playground:  http://localhost:8080/playground")
	log.Fatal(http.ListenAndServe(conf.Port, nil))
}

func startPostgres(connectingString string) interfaces.StoreIface {
	pool := postgres.New(connectingString)
	return databases.NewPostgresRepository(pool.Pool)
}

func startInMemory() interfaces.StoreIface {
	return databases.NewInMemRepository()
}
