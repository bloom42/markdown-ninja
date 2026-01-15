package service

import (
	"context"
	"time"

	"github.com/skerkour/stdx-go/retry"
	"markdown.ninja/pkg/services/content"
)

func (service *ContentService) JobDeleteAssetsDataWithPrefix(ctx context.Context, input content.JobDeleteAssetsDataWithPrefix) (err error) {
	err = retry.Do(func() error {
		return service.storage.DeleteObjectsWithPrefix(ctx, input.Prefix)
	}, retry.Context(ctx), retry.Attempts(5), retry.Delay(1*time.Minute))
	return
}
