package service

import (
	"context"
	"time"

	"github.com/skerkour/stdx-go/opt"
	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/contacts"
	"markdown.ninja/pkg/services/events"
	"markdown.ninja/pkg/services/site"
)

func (service *SiteService) Unsubscribe(ctx context.Context, input site.UnsubscribeInput) (err error) {
	actor := service.contactsService.CurrentContact(ctx)

	service.kernel.SleepAuth()

	httpCtx := httpctx.FromCtx(ctx)
	now := time.Now().UTC()

	website, err := service.websitesService.FindWebsiteByDomain(ctx, service.db, httpCtx.Hostname)
	if err != nil {
		return
	}

	contactID, err := service.contactsService.ParseAndVerifyUnsubscribeToken(input.Token)
	if err != nil {
		return
	}

	if actor != nil && !actor.ID.Equal(contactID) {
		return contacts.ErrUnsubscribeLinkIsNotValid
	}

	contact, err := service.contactsService.FindContact(ctx, service.db, contactID)
	if err != nil {
		return
	}

	if !website.ID.Equal(contact.WebsiteID) {
		return contacts.ErrUnsubscribeLinkIsNotValid
	}

	if contact.Email != input.Email {
		return contacts.ErrUnsubscribeLinkIsNotValid
	}

	if contact.SubscribedToNewsletterAt == nil {
		return
	}

	contact.UpdatedAt = now
	contact.SubscribedToNewsletterAt = nil

	updateContactInput := contacts.UpdateContactInput{
		ID:                     contact.ID,
		SubscribedToNewsletter: opt.Bool(false),
	}
	err = service.contactsService.UpdateContactInternal(ctx, service.db, &contact, updateContactInput)
	if err != nil {
		return
	}

	trackEventInput := events.TrackUnsubscribedFromNewsletterInput{
		WebsiteID: website.ID,
	}
	service.eventsService.TrackUnsubscribedFromNewsletter(ctx, trackEventInput)

	return
}
