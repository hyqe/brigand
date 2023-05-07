package main

import (
	"context"

	"github.com/hyqe/brigand/internal/server"
)

func main() { server.Run(context.Background()) }
