package server

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/hyqe/timber"
)

func Test_graceful(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jack := timber.NewJack()
	srv := &http.Server{}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		graceful(ctx, jack, srv)
	}()
	cancel()
	wg.Wait()
}
