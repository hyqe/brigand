package server

import (
	"context"
	"net/http"
	"sync"
	"testing"
)

func Test_graceful(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	srv := &http.Server{}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		graceful(ctx, srv)
	}()
	cancel()
	wg.Wait()
}
