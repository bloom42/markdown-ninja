package service

import (
	"context"
	"encoding/csv"
	"strings"
	"time"

	"github.com/bloom42/stdx-go/countries"
	"github.com/bloom42/stdx-go/db"
	"github.com/bloom42/stdx-go/guid"
	"github.com/bloom42/stdx-go/log/slogx"
	"markdown.ninja/pkg/errs"
	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/contacts"
	"markdown.ninja/pkg/services/events"
)

func (service *ContactsService) ImportContacts(ctx context.Context, input contacts.ImportContactsInput) (ret []contacts.Contact, err error) {
	logger := slogx.FromCtx(ctx)
	httpCtx := httpctx.FromCtx(ctx)

	_, err = service.kernel.CurrentUserID(ctx)
	if err != nil {
		return
	}
	if !httpCtx.AccessToken.IsAdmin {
		return ret, errs.PermissionDenied("Please contact support to import contacts")
	}

	// err = service.websitesService.CheckUserIsStaff(ctx, service.db, actorID, input.WebsiteID)
	// if err != nil {
	// 	return
	// }

	csvReader := csv.NewReader(strings.NewReader(strings.TrimSpace(input.ContactsCsv)))
	csvRecords, err := csvReader.ReadAll()
	if err != nil {
		logger.Warn("contacts.ImportContacts: parsing CSV", slogx.Err(err))
		return ret, contacts.ErrImportingContacts
	}

	if len(csvRecords) == 0 {
		return []contacts.Contact{}, nil
	}

	if len(csvRecords[0]) != 3 {
		return ret, contacts.ErrImportCsvHeaderisNotValid
	}

	headerEmail := strings.ToLower(strings.TrimSpace(csvRecords[0][0]))
	headerName := strings.ToLower(strings.TrimSpace(csvRecords[0][1]))
	headerSubscribedAt := strings.ToLower(strings.TrimSpace(csvRecords[0][2]))
	if headerEmail != "email" || headerName != "name" || headerSubscribedAt != "subscribed_at" {
		return ret, contacts.ErrImportCsvHeaderisNotValid
	}

	now := time.Now().UTC()
	importedContacts := make([]contacts.Contact, 0, len(csvRecords))

	// we start at 1 because row 0 is for the CSV header
	for i := 1; i < len(csvRecords); i += 1 {
		if len(csvRecords[i]) != 3 {
			err = contacts.ErrImportingContacts
			return
		}

		email := strings.ToLower(strings.TrimSpace(csvRecords[i][0]))

		err = service.ValidateContactEmail(ctx, email, false)
		if err != nil {
			return
		}

		name := strings.TrimSpace(csvRecords[i][1])
		if name == "" {
			name = service.extractNameFromEmail(email)
		}

		err = service.ValidateContactName(name)
		if err != nil {
			return
		}

		subscribedToNewsletterAtStr := strings.ToUpper(strings.TrimSpace(csvRecords[i][2]))
		var subscribedToNewsletterAt *time.Time
		if subscribedToNewsletterAtStr != "" {
			var subscribedAtTmp time.Time
			subscribedAtTmp, err = time.Parse(time.RFC3339, subscribedToNewsletterAtStr)
			if err != nil {
				err = contacts.ErrImportSubscribedAtIsNotValid
				return
			}
			subscribedToNewsletterAt = &subscribedAtTmp
		}

		importedContact := contacts.Contact{
			ID:                           guid.NewTimeBased(),
			CreatedAt:                    now,
			UpdatedAt:                    now,
			Name:                         name,
			Email:                        email,
			SubscribedToNewsletterAt:     subscribedToNewsletterAt,
			SubscribedToProductUpdatesAt: &now,
			Verified:                     true,
			Country:                      countries.CodeUnknown,
			FailedSignupAttempts:         0,
			SignupCodeHash:               "",
			StripeCustomerID:             nil,
			WebsiteID:                    input.WebsiteID,
		}
		importedContacts = append(importedContacts, importedContact)
	}

	// we first buffer events to not save save them if an error happens during the transaction
	eventsToSave := make([]events.TrackSubscribedToNewsletterInput, 0, len(importedContacts))

	err = service.db.Transaction(ctx, func(tx db.Tx) (txErr error) {

		// TODO: improve perfs
		for _, importedContact := range importedContacts {
			// if contacts already exists but is not verified yet, we mark it as verified
			// and update the relevant information
			var existingContact contacts.Contact
			existingContact, txErr = service.repo.FindContactByEmail(ctx, tx, input.WebsiteID, importedContact.Email)
			if txErr != nil {
				if !errs.IsNotFound(txErr) {
					return txErr
				}

				txErr = service.repo.CreateContact(ctx, tx, importedContact)
				if txErr != nil {
					return txErr
				}

				if importedContact.SubscribedToNewsletterAt != nil {
					trackEventInput := events.TrackSubscribedToNewsletterInput{
						WebsiteID: input.WebsiteID,
					}
					eventsToSave = append(eventsToSave, trackEventInput)
				}

			} else {
				// if contact exists but is not verified or not subscribed
				if !existingContact.Verified ||
					existingContact.SubscribedToNewsletterAt != importedContact.SubscribedToNewsletterAt {
					existingContact.Verified = true
					existingContact.FailedSignupAttempts = 0
					existingContact.SubscribedToNewsletterAt = importedContact.SubscribedToNewsletterAt
					txErr = service.repo.UpdateContact(ctx, tx, existingContact)
					if txErr != nil {
						return txErr
					}

					if existingContact.SubscribedToNewsletterAt != nil {
						trackEventInput := events.TrackSubscribedToNewsletterInput{
							WebsiteID: input.WebsiteID,
						}
						eventsToSave = append(eventsToSave, trackEventInput)
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return
	}

	for _, event := range eventsToSave {
		service.eventsService.TrackSubscribedToNewsletter(ctx, event)
	}

	ret = importedContacts

	return
}
