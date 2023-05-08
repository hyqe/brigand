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
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

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

	httpServer := http.Server{
		Addr:    cfg.Addr(),
		Handler: middleware(routes),
	}

	// start the http server
	go func() {
		jack.Alert(fmt.Sprintf("listening on '%v'\n", httpServer.Addr))
		err = httpServer.ListenAndServe()
		if err != nil {
			jack.Error(err)
		}
	}()

	// wait for context to cancel
	shutdown(ctx, func() {
		jack.Alert("shutting down http server")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		httpServer.Shutdown(ctx)
	})
}

// shutdown listens for context cancels and then calls the drain functions.
// it is up to the individual functions to decide how to shutdown, and if
// there should be a max timeout before forcefully shutting down.
func shutdown(ctx context.Context, fns ...func()) {
	<-ctx.Done()

	// shutdown all funcs at once, wait for them to finish.
	var wg sync.WaitGroup
	for _, fn := range fns {
		wg.Add(1)
		go func(fn func()) {
			defer func() {
				if v := recover(); v != nil {
					fmt.Println(v)
				}
			}()
			defer wg.Done()
			fn()
		}(fn)
	}
	wg.Wait()
}
