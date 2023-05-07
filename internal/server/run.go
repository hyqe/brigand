package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hyqe/brigand/internal/storage"
	"github.com/hyqe/timber"
)

func Run(ctx context.Context) {
	cfg, err := GetConfig()
	if err != nil {
		panic(err)
	}

	mongoClient, err := storage.NewMongoClient(ctx, cfg.MongoUri)
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(ctx)

	jack := timber.NewJack(
		timber.WithLevel(cfg.Level),
	)

	routes := Routes(
		storage.NewMongoMetadataClient(mongoClient),
	)

	middleware := timber.NewMiddleware(jack)

	fmt.Printf("listening on '%v'\n", cfg.Addr())
	http.ListenAndServe(cfg.Addr(), middleware(routes))
}
