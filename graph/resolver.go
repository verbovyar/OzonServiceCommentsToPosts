package graph

import (
	"ozonProject/internal/pubsub"
	"ozonProject/internal/service"
)

type Resolver struct {
	Service *service.Service
	Bus     *pubsub.Bus
}
