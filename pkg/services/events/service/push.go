package service

import (
	"context"
	"log/slog"

	"github.com/skerkour/stdx-go/log/slogx"
	"markdown.ninja/pkg/services/events"
)

func (service *Service) Push(ctx context.Context, event events.Event) {
	logger := slogx.FromCtx(ctx)
	if event.Data == nil {
		logger.Error("events.Push: data is null", slog.String("event.type", event.Type.String()))
		return
	}

	service.eventsBuffer.Push(event)
}
