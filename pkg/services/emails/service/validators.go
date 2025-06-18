package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"markdown.ninja/pkg/services/emails"
	"markdown.ninja/pkg/services/kernel"
)

func (service *EmailsService) validateNewsletterSubject(subject string) (err error) {
	if len(subject) < emails.NewsletterSubjectMinSize {
		return emails.ErrNewsletterSubjectIsTooShort
	}

	if len(subject) > emails.NewsletterSubjectMaxSize {
		return emails.ErrNewsletterSubjectIsTooLong
	}

	if !utf8.ValidString(subject) {
		return emails.ErrNewsletterSubjectIsNotValid
	}

	return nil
}

func (service *EmailsService) validateNewsletterBodyMarkdown(input string) (err error) {
	if len(input) > emails.NewsletterContentMarkdownMaxSize {
		return emails.ErrNewsletterBodyIsTooLarge
	}

	if !utf8.ValidString(input) {
		return emails.ErrNewsletterBodyIsNotValid
	}

	return nil
}

func (service *EmailsService) validateNewsletterScheduledFor(scheduledFor *time.Time) (err error) {
	if scheduledFor == nil {
		return nil
	}

	now := time.Now().UTC()
	if scheduledFor.Before(now) {
		return emails.ErrNewsletterScheduledForIsInThePast
	}

	return nil
}

func (service *EmailsService) validateSenderEmailAddress(ctx context.Context, email string) (err error) {
	err = service.kernel.ValidateEmail(ctx, email, true)
	if err != nil {
		return err
	}

	// No need to add more check (out of bound...) because we already verified that email has 1 and only 1
	// "@" character with the previous call to ValidateEmail
	domain := strings.Split(email, "@")[1]

	// check that it's not the domain used by Markdown Ninja
	if domain == service.httpConfig.WebappDomain || domain == service.httpConfig.WebsitesRootDomain {
		return kernel.ErrEmailIsNotValid
	}

	// if len(mxRecords) == 0 {
	// 	return errs.InvalidArgument("MX records are required for the email sending domain")
	// }

	return nil
}
