package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/store"
)

func (service *StoreService) FindCompletedOrdersForContact(ctx context.Context, db db.Queryer, contactID guid.GUID) (orders []store.Order, err error) {
	orders, err = service.repo.FindOrdersWithStatusForContact(ctx, db, contactID, store.OrderStatusCompleted)
	return
}
