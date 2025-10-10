package main

import (
	"fmt"
	"log"
	"net/http"
	"ozonProject/config"
	"ozonProject/graph"
	"ozonProject/internal/pubsub"
	"ozonProject/internal/service"
	"ozonProject/internal/storage"

	"ozonProject/pkg/postgres"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

const (
	queryPath      = "/query"
	playgroundPath = "/playground"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	runApp(config)
}

func usePostgres(config config.Config) storage.Storage {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	pool := postgres.New(connectionString)
	return storage.NewPostgresStorage(pool.Pool)
}

func useInMemory() storage.Storage {
	return storage.NewInMemoryStorage()
}

func runApp(config config.Config) {
	var repo storage.Storage
	if config.PersistanceEnabled {
		repo = usePostgres(config)
	} else {
		repo = useInMemory()
	}

	service := service.New(repo)
	bus := pubsub.New()

	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{Service: service, Bus: bus},
	}))
	server.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	http.Handle(playgroundPath, playground.Handler("Playground", queryPath))
	http.Handle(queryPath, server)

	log.Printf("listening on %s", config.AppPort)
	log.Printf("Sandbox:  http://localhost:%s%s", config.AppPort, playgroundPath)
	log.Fatal(http.ListenAndServe(config.AppPort, nil))
}
