package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

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

	graceful(ctx, jack, &http.Server{
		Addr:    cfg.Addr(),
		Handler: middleware(routes),
	})
}

func graceful(
	ctx context.Context,
	jack timber.Jack,
	srv *http.Server,
) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	jack.Alert(fmt.Sprintf("listening on '%v'\n", srv.Addr))
	go srv.ListenAndServe()

	<-ctx.Done()
	jack.Alert("shutting down http server")
	srv.Shutdown(context.Background())
}
