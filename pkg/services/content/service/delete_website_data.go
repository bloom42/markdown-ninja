package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"github.com/skerkour/stdx-go/log/slogx"
	"github.com/skerkour/stdx-go/opt"
	"github.com/skerkour/stdx-go/queue"
	"markdown.ninja/pkg/errs"
	"markdown.ninja/pkg/services/content"
)

func (service *ContentService) DeleteWebsiteData(ctx context.Context, tx db.Queryer, websiteID guid.GUID) (err error) {
	logger := slogx.FromCtx(ctx)

	job := queue.NewJobInput{
		Data: content.JobDeleteAssetsDataWithPrefix{
			Prefix: service.getStoragePrefixForWebsite(websiteID),
		},
		Timeout:    opt.Ptr(int64(600)),
		RetryDelay: opt.Ptr(int64(3600)),
	}
	err = service.queue.Push(ctx, tx, job)
	if err != nil {
		errMessage := "content.DeleteWebsiteAssets: Pushing DeleteAssetsWithPrefix job to queue"
		logger.Error(errMessage, slogx.Err(err))
		err = errs.Internal(errMessage, err)
		return
	}

	err = service.repo.DeleteWebsiteAssets(ctx, tx, websiteID)
	if err != nil {
		return
	}

	return
}
