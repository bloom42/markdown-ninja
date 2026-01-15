package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/content"
)

func (service *ContentService) FindWebsiteAssetByID(ctx context.Context, db db.Queryer, websiteID, assetID guid.GUID) (asset content.Asset, err error) {
	asset, err = service.repo.FindAssetByID(ctx, db, assetID)
	if err != nil {
		return
	}
	if !asset.WebsiteID.Equal(websiteID) {
		err = content.ErrAssetNotFound
		return
	}
	return
}
