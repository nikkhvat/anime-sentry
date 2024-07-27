package handlers

import "context"

type Bot interface {
	Start(ctx context.Context) error
}
