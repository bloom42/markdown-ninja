package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"markdown.ninja/pkg/services/store"
)

func (service *StoreService) hydrateProduct(ctx context.Context, db db.Queryer, product *store.Product) (err error) {
	product.Content, err = service.repo.FindProductPagesForProduct(ctx, db, product.ID)
	if err != nil {
		return err
	}

	product.Assets, err = service.contentService.FindProductAssets(ctx, db, product.ID)
	if err != nil {
		return err
	}

	return nil
}
