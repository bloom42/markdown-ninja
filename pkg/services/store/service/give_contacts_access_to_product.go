package service

import (
	"context"
	"strings"
	"time"

	"github.com/bloom42/stdx-go/countries"
	"github.com/bloom42/stdx-go/db"
	"github.com/bloom42/stdx-go/opt"
	"github.com/bloom42/stdx-go/slicesx"
	"markdown.ninja/pkg/errs"
	"markdown.ninja/pkg/services/contacts"
	"markdown.ninja/pkg/services/store"
)

func (service *StoreService) GiveContactsAccessToProduct(ctx context.Context, input store.GiveContactsAccessToProductInput) (err error) {
	actorID, err := service.kernel.CurrentUserID(ctx)
	if err != nil {
		return
	}

	product, err := service.repo.FindProductByID(ctx, service.db, input.ProductID)
	if err != nil {
		return
	}

	err = service.websitesService.CheckUserIsStaff(ctx, service.db, actorID, product.WebsiteID)
	if err != nil {
		return
	}

	emails := slicesx.Unique(input.Emails)

	now := time.Now().UTC()

	err = service.db.Transaction(ctx, func(tx db.Tx) (txErr error) {
		for _, email := range emails {
			email = strings.ToLower(strings.TrimSpace(email))
			if email == "" {
				continue
			}

			txErr = service.contactsService.ValidateContactEmail(ctx, email, false)
			if txErr != nil {
				return txErr
			}

			var contact contacts.Contact
			contact, txErr = service.contactsService.FindContactByEmail(ctx, tx, product.WebsiteID, email)
			if txErr != nil {
				if !errs.IsNotFound(txErr) {
					return txErr
				}

				// if contact is not found we create it
				createContactInput := contacts.CreateContactInternalInput{
					Email:     email,
					Name:      "",
					Verified:  true,
					Country:   countries.CodeUnknown,
					WebsiteID: product.WebsiteID,
				}
				contact, txErr = service.contactsService.CreateContactInternal(ctx, tx, createContactInput)
				if txErr != nil {
					return txErr
				}
			} else {
				// if contact is found but not verified, mark as verified
				if !contact.Verified {
					updateContactInput := contacts.UpdateContactInput{
						ID:       contact.ID,
						Verified: opt.Bool(true),
					}
					txErr = service.contactsService.UpdateContactInternal(ctx, tx, &contact, updateContactInput)
					if txErr != nil {
						return txErr
					}
				}
			}

			_, txErr = service.repo.FindContactProductAccess(ctx, tx, contact.ID, product.ID)
			if txErr == nil {
				// if contact already has access to product we do nothing
				continue
			} else {
				if !errs.IsNotFound(txErr) {
					return txErr
				}
				txErr = nil
			}

			contactProductAccess := store.ContactProductAccess{
				CreatedAt: now,
				ContactID: contact.ID,
				ProductID: product.ID,
			}
			txErr = service.repo.CreateContactProductAccess(ctx, tx, contactProductAccess)
			if txErr != nil {
				return txErr
			}
		}

		return nil
	})
	if err != nil {
		return
	}

	return
}
