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

	timber.Default.Apply(timber.WithLevel(cfg.Level))

	mongoClient, err := storage.NewMongoClient(ctx, cfg.MongoUri)
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(ctx)

	s3Sess, err := storage.NewS3Session(cfg.DOSpacesRegion, cfg.DOSpacesEndpoint, cfg.DOSpacesAccessKey, cfg.DOSpacesSecretKey)
	if err != nil {
		panic(err)
	}

	routes := Routes(
		storage.NewMongoMetadataClient(mongoClient),
		storage.NewS3FileDownloader(s3Sess, cfg.DOSpacesBucket),
		storage.NewS3FileUploader(s3Sess, cfg.DOSpacesBucket),
		cfg.SymlinkSecret,
	)

	log := timber.NewMiddleware()

	sudo := SudoMiddlware(cfg.Sudo)

	graceful(ctx, &http.Server{
		Addr:    cfg.Addr(),
		Handler: log(sudo(routes)),
	})
}

func graceful(
	ctx context.Context,
	srv *http.Server,
) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		timber.Alert(fmt.Sprintf("listening on '%v'\n", srv.Addr))

		if err := srv.ListenAndServe(); err != nil {
			timber.Error(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		timber.Alert("shutting down http server")

		// if the server takes more then a minute to shutdown, something is seriously wrong.
		// A minute is overkill, but we just need some failsafe that will ensure the process
		// is killed eventually.
		timeout, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		if err := srv.Shutdown(timeout); err != nil {
			timber.Error(err)
		}
	}()

	wg.Wait()
}
