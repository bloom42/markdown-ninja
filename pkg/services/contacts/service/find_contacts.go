package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/contacts"
)

func (service *ContactsService) FindContactsByEmail(ctx context.Context, db db.Queryer, websiteID guid.GUID, emails []string) (contacts []contacts.Contact, err error) {
	contacts, err = service.repo.FindContactsByEmails(ctx, db, websiteID, emails)
	return
}
