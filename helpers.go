package main

import (
	"context"
	"time"
)

func getContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return context.WithTimeout(ctx, 3*time.Second)
}
