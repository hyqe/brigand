package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

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

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		jack.Alert(fmt.Sprintf("listening on '%v'\n", srv.Addr))

		if err := srv.ListenAndServe(); err != nil {
			jack.Error(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		jack.Alert("shutting down http server")

		// if the server takes more then a minute to shutdown, something is seriously wrong.
		// A minute is overkill, but we just need some failsafe that will ensure the process
		// is killed eventually.
		timeout, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		if err := srv.Shutdown(timeout); err != nil {
			jack.Error(err)
		}
	}()

	wg.Wait()
}
