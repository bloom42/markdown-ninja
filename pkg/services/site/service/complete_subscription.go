package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/bloom42/stdx-go/countries"
	"github.com/bloom42/stdx-go/crypto"
	"github.com/bloom42/stdx-go/db"
	"github.com/bloom42/stdx-go/opt"
	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/contacts"
	"markdown.ninja/pkg/services/events"
	"markdown.ninja/pkg/services/kernel"
	"markdown.ninja/pkg/services/site"
)

// TODO: event
func (service *SiteService) CompleteSubscription(ctx context.Context, input site.CompleteSubscriptionInput) (retContact site.Contact, err error) {
	authenticatedContact := service.contactsService.CurrentContact(ctx)
	if authenticatedContact != nil {
		err = kernel.ErrMustNotBeAuthenticated
		return
	}

	service.kernel.SleepAuth()
	httpCtx := httpctx.FromCtx(ctx)
	contactID := input.ContactID
	code := strings.ToLower(strings.TrimSpace(input.Code))
	domain := httpCtx.Hostname
	now := time.Now().UTC()

	website, err := service.websitesService.FindWebsiteByDomain(ctx, service.db, domain)
	if err != nil {
		return
	}

	contact, err := service.contactsService.FindContact(ctx, service.db, contactID)
	if err != nil {
		return
	}

	if contact.Verified {
		err = site.ErrAccountAlreadyExists
		return
	}

	if !contact.WebsiteID.Equal(website.ID) {
		err = site.ErrAccountNotFound
		service.kernel.SleepAuthFailure()
		return
	}

	if contact.FailedSignupAttempts >= site.SignupMaxAttempts {
		err = site.ErrAuthCodeInvalidOrExpired
		service.kernel.SleepAuthFailure()
		return
	}

	since := now.Sub(contact.UpdatedAt)
	if since >= time.Hour {
		err = site.ErrAuthCodeInvalidOrExpired
		service.kernel.SleepAuthFailure()
		return
	}

	if !crypto.VerifyPasswordHash([]byte(code), contact.SignupCodeHash) {
		err = kernel.ErrAuthCodeIsNotValid
		// session.UpdatedAt = now
		updateContactInput := contacts.UpdateContactInput{
			ID:                   contact.ID,
			FailedSignupAttempts: opt.Int64(contact.FailedSignupAttempts + 1),
		}
		err = service.contactsService.UpdateContactInternal(ctx, service.db, &contact, updateContactInput)
		if err != nil {
			return
		}
		err = site.ErrAuthCodeInvalidOrExpired
		service.kernel.SleepAuthFailure()
		return
	}

	country := httpCtx.Client.CountryCode
	if country == countries.CodeUnknown {
		country = contact.Country
	}

	var sessionCookie *http.Cookie
	err = service.db.Transaction(ctx, func(tx db.Tx) (txErr error) {
		updateContactInput := contacts.UpdateContactInput{
			ID:                     contact.ID,
			SubscribedToNewsletter: opt.Bool(true),
			Country:                &country,
			FailedSignupAttempts:   opt.Int64(0),
			SignupCodeHash:         opt.String(""),
			Verified:               opt.Bool(true),
		}
		txErr = service.contactsService.UpdateContactInternal(ctx, tx, &contact, updateContactInput)
		if txErr != nil {
			return txErr
		}

		createSessionInput := contacts.CreateSessionInput{
			Verified:  true,
			ContactID: contact.ID,
			WebsiteID: website.ID,
		}
		_, sessionCookie, txErr = service.contactsService.CreateSession(ctx, tx, createSessionInput)
		return txErr
	})
	if err != nil {
		return
	}

	if sessionCookie != nil {
		httpCtx.Response.Cookies = append(httpCtx.Response.Cookies, *sessionCookie)
	}

	retContact = service.convertContact(contact)

	trackEventInput := events.TrackSubscribedToNewsletterInput{
		WebsiteID: website.ID,
	}
	service.eventsService.TrackSubscribedToNewsletter(ctx, trackEventInput)

	return
}
