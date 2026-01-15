package service

import (
	"context"
	"time"

	"github.com/skerkour/stdx-go/retry"
	"markdown.ninja/pkg/services/content"
)

func (service *ContentService) JobDeleteAssetData(ctx context.Context, input content.JobDeleteAssetData) (err error) {
	err = retry.Do(func() error {
		return service.storage.DeleteObject(ctx, input.StorageKey)
	}, retry.Context(ctx), retry.Attempts(5), retry.Delay(1*time.Minute))

	return
}
