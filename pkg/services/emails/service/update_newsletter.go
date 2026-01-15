package service

import (
	"context"
	"strings"
	"time"

	"github.com/skerkour/stdx-go/crypto/blake3"
	"markdown.ninja/pkg/services/emails"
)

// TODO: update website's used storage? See also CreateNewsletter and DeleteNewsletter
func (service *EmailsService) UpdateNewsletter(ctx context.Context, input emails.UpdateNewsletterInput) (newsletter emails.Newsletter, err error) {
	actorID, err := service.kernel.CurrentUserID(ctx)
	if err != nil {
		return
	}

	newsletter, err = service.repo.FindNewsletterByID(ctx, service.db, input.ID)
	if err != nil {
		return
	}

	err = service.websitesService.CheckUserIsStaff(ctx, service.db, actorID, newsletter.WebsiteID)
	if err != nil {
		return
	}

	now := time.Now().UTC()
	newsletter.UpdatedAt = now
	newsletter.Subject = strings.TrimSpace(input.Subject)
	newsletter.ScheduledFor = input.ScheduledFor

	if input.BodyMarkdown != nil {
		newsletter.BodyMarkdown = *input.BodyMarkdown
		err = service.validateNewsletterBodyMarkdown(newsletter.BodyMarkdown)
		if err != nil {
			return
		}

		newsletter.Size = int64(len(newsletter.BodyMarkdown))
		bodyHash := blake3.Sum256([]byte(newsletter.BodyMarkdown))
		newsletter.Hash = bodyHash[:]
	}

	if newsletter.ScheduledFor != nil {
		scheduledForUtc := newsletter.ScheduledFor.UTC()
		newsletter.ScheduledFor = &scheduledForUtc
		err = service.validateNewsletterScheduledFor(newsletter.ScheduledFor)
		if err != nil {
			return
		}
	}

	err = service.validateNewsletterSubject(newsletter.Subject)
	if err != nil {
		return
	}

	err = service.repo.UpdateNewsletter(ctx, service.db, newsletter)
	if err != nil {
		return
	}

	return
}
