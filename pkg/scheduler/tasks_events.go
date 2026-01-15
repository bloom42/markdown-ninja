package scheduler

import (
	"context"

	"github.com/skerkour/stdx-go/log/slogx"
	"github.com/skerkour/stdx-go/queue"
	"markdown.ninja/pkg/services/events"
)

func (scheduler *Scheduler) eventsDispatchRotateAnonymousIDSalt(ctx context.Context) {
	job := queue.NewJobInput{
		Data: events.JobRotateAnonymousIDSalt{},
	}
	err := scheduler.queue.Push(ctx, nil, job)
	if err != nil {
		logger := slogx.FromCtx(ctx)
		logger.Error("scheduler.DispatchRotateAnonymousIDSalt: Pushing job to queue", slogx.Err(err))
		return
	}
}
