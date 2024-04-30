package handlers

import "context"

type EventHandler interface {
	Handle(ctx context.Context)
}
