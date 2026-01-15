package service

import (
	"context"

	"github.com/skerkour/stdx-go/queue"
	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/kernel"
)

func (service *KernelService) ListFailedBackgroundJobs(ctx context.Context, input kernel.EmptyInput) (ret kernel.PaginatedResult[queue.Job], err error) {
	httpCtx := httpctx.FromCtx(ctx)

	_, err = service.CurrentUserID(ctx)
	if err != nil {
		return
	}

	if !httpCtx.AccessToken.IsAdmin {
		service.SleepAuth()
		err = kernel.ErrPermissionDenied
		return
	}

	ret.Data, err = service.queue.GetFailedJobs(ctx)
	return
}
