package scheduler

import (
	"context"

	"github.com/bloom42/stdx-go/log/slogx"
	"github.com/bloom42/stdx-go/opt"
	"github.com/bloom42/stdx-go/queue"
	"markdown.ninja/pkg/services/organizations"
)

func (scheduler *Scheduler) organizationsDispatchSendUsageData(ctx context.Context) {
	job := queue.NewJobInput{
		Data:       organizations.JobDispatchSendUsageData{},
		Timeout:    opt.Ptr(int64(300)),
		RetryDelay: opt.Ptr(int64(300)),
	}
	err := scheduler.queue.Push(ctx, nil, job)
	if err != nil {
		logger := slogx.FromCtx(ctx)
		logger.Error("scheduler.DispatchSendUsageData: Pushing job to queue", slogx.Err(err))
		return
	}
}

func (scheduler *Scheduler) organizationsDispatchInvoiceMonthlyUsage(ctx context.Context) {
	job := queue.NewJobInput{
		Data:       organizations.JobDispatchInvoiceMonthlyUsage{},
		Timeout:    opt.Ptr(int64(300)),
		RetryDelay: opt.Ptr(int64(300)),
	}
	err := scheduler.queue.Push(ctx, nil, job)
	if err != nil {
		logger := slogx.FromCtx(ctx)
		logger.Error("scheduler.DispatchInvoiceMonthlyUsage: Pushing job to queue", slogx.Err(err))
		return
	}
}
