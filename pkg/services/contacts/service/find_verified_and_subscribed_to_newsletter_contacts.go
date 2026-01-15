package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/contacts"
)

func (service *ContactsService) FindVerifiedAndSubscribedToNewsletterContacts(ctx context.Context, db db.Queryer, websiteID guid.GUID) (contacts []contacts.Contact, err error) {
	contacts, err = service.repo.FindVerifiedAndSubscribedToNewsletterContacts(ctx, db, websiteID)
	return
}
