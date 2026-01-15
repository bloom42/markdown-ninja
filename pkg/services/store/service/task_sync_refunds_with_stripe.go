package service

import (
	"context"

	"github.com/skerkour/stdx-go/log/slogx"
	"github.com/skerkour/stdx-go/queue"
	"markdown.ninja/pkg/services/store"
)

// TODO: batch push jobs
func (service *StoreService) TaskSyncRefundsWithStripe(ctx context.Context) {
	logger := slogx.FromCtx(ctx)

	pendingRefunds, err := service.repo.FindPendingRefunds(ctx, service.db)
	if err != nil {
		logger.Error("stroe.TaskSyncRefundsFromStripe: error finding pending refunds", slogx.Err(err))
		return
	}

	for _, refund := range pendingRefunds {
		job := queue.NewJobInput{
			Data: store.JobSyncRefundWithStripe{
				RefundID: refund.ID,
			},
		}
		err := service.queue.Push(ctx, nil, job)
		if err != nil {
			logger.Error("stroe.TaskSyncRefundsFromStripe: Pushing job to queue", slogx.Err(err))
			continue
		}
	}
}
