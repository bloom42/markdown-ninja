package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/store"
)

func (service *StoreService) FindProductsForContact(ctx context.Context, db db.Queryer, contactID guid.GUID) (products []store.Product, err error) {
	products, err = service.repo.FindProductsForContact(ctx, db, contactID)
	return
}
