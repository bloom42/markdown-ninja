package service

import (
	"context"

	"github.com/skerkour/stdx-go/log/slogx"
	"github.com/skerkour/stdx-go/queue"
	"markdown.ninja/pkg/services/contacts"
)

func (service *ContactsService) TaskDeleteOldUnverifiedSessions(ctx context.Context) {
	logger := slogx.FromCtx(ctx)

	job := queue.NewJobInput{
		Data: contacts.JobDeleteOldUnverifiedSessions{},
	}
	err := service.queue.Push(ctx, nil, job)
	if err != nil {
		logger.Error("contacts.TaskDeleteOldUnverifiedSessions: Pushing job to queue", slogx.Err(err))
		return
	}
}
