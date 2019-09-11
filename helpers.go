package main

import (
	"context"
	"time"
)

// Errors and string checks
const (
	ErrCodeNotFound = 404
	ErrNotFound     = "The requested resource was not found"
)

func getContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return context.WithTimeout(ctx, 8*time.Second)
}
