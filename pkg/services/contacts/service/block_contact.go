package service

import (
	"context"
	"time"

	"github.com/bloom42/stdx-go/db"
	"markdown.ninja/pkg/services/contacts"
)

func (service *ContactsService) BlockContact(ctx context.Context, input contacts.BlockContactInput) (contact contacts.Contact, err error) {
	actorID, err := service.kernel.CurrentUserID(ctx)
	if err != nil {
		return
	}

	contact, err = service.repo.FindContactByID(ctx, service.db, input.ID)
	if err != nil {
		return
	}

	err = service.websitesService.CheckUserIsStaff(ctx, service.db, actorID, contact.WebsiteID)
	if err != nil {
		return
	}

	if contact.BlockedAt != nil {
		err = contacts.ErrContactIsAlreadyBlocked
		return
	}

	now := time.Now().UTC()
	contact.UpdatedAt = now
	contact.BlockedAt = &now
	err = service.db.Transaction(ctx, func(tx db.Tx) (txErr error) {
		txErr = service.repo.UpdateContact(ctx, tx, contact)
		if txErr != nil {
			return
		}

		txErr = service.repo.DeleteSessionsForContact(ctx, tx, contact.ID)
		if txErr != nil {
			return
		}

		return nil
	})
	if err != nil {
		return
	}

	return
}
