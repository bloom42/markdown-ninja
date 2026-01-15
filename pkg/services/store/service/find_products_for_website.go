package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/store"
)

func (service *StoreService) FindProductsForWebsite(ctx context.Context, db db.Queryer, websiteID guid.GUID, limit int64) (products []store.Product, err error) {
	products, err = service.repo.FindProductsByWebsiteID(ctx, db, websiteID, limit)
	return
}
