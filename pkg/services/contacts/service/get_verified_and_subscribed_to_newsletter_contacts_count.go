package service

import (
	"context"

	"github.com/bloom42/stdx-go/db"
	"github.com/bloom42/stdx-go/guid"
)

func (service *ContactsService) GetVerifiedAndSubscribedToNewsletterContactsCount(ctx context.Context, db db.Queryer, websiteID guid.GUID) (count int64, err error) {
	count, err = service.repo.GetVerifiedAndSubscribedToNewsletterContactsCount(ctx, service.db, websiteID)
	return
}
