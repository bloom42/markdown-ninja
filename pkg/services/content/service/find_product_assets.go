package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/content"
)

func (service *ContentService) FindProductAssets(ctx context.Context, db db.Queryer, productID guid.GUID) (assets []content.Asset, err error) {
	assets, err = service.repo.FindProductAssets(ctx, db, productID)
	return
}
