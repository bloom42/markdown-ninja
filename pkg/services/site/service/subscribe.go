package service

import (
	"context"
	"strings"
	"time"

	"github.com/skerkour/stdx-go/countries"
	"github.com/skerkour/stdx-go/crypto"
	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/log/slogx"
	"github.com/skerkour/stdx-go/opt"
	"github.com/skerkour/stdx-go/queue"
	"github.com/skerkour/stdx-go/randutil"
	"markdown.ninja/pkg/errs"
	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/contacts"
	"markdown.ninja/pkg/services/site"
)

func (service *SiteService) Subscribe(ctx context.Context, input site.SubscribeInput) (ret site.SubscribeOutput, err error) {

	// sleep to prevent spam and bruteforce
	service.kernel.SleepAuth()

	httpCtx := httpctx.FromCtx(ctx)
	hostname := httpCtx.Hostname
	email := strings.ToLower(strings.TrimSpace(input.Email))
	name := strings.ToLower(strings.TrimSpace(input.Name))
	unverifiedContactAlreadyExists := false
	logger := slogx.FromCtx(ctx)

	err = service.kernel.ValidateEmail(ctx, email, true)
	if err != nil {
		return
	}

	err = service.contactsService.ValidateContactName(name)
	if err != nil {
		return
	}

	website, err := service.websitesService.FindWebsiteByDomain(ctx, service.db, hostname)
	if err != nil {
		return
	}

	if !service.rateLimiter.IsAllowed("SiteService.Subscribe", website.ID.Bytes(), httpCtx.Client.IP.AsSlice(), time.Hour, 30) {
		err = errs.TooManyRequests()
		return
	}

	contact, err := service.contactsService.FindContactByEmail(ctx, service.db, website.ID, email)
	if err == nil {
		if contact.Verified {
			err = site.ErrAccountAlreadyExists
			service.kernel.SleepAuthFailure()
			return
		}
		unverifiedContactAlreadyExists = true
	} else {
		if !errs.IsNotFound(err) {
			return
		}
		err = nil
	}

	randomGenerator := crypto.NewRandomGenerator()
	codeBytes := randutil.RandAlphabet(randomGenerator, []byte(site.AuthCodeAlphabet), site.AuthCodeLength)
	code := string(codeBytes)
	codeHash := crypto.HashPassword(codeBytes, site.AuthCodeHashParams)

	country := httpCtx.Client.CountryCode
	if country == countries.CodeUnknown && unverifiedContactAlreadyExists {
		country = contact.Country
	}

	err = service.db.Transaction(ctx, func(tx db.Tx) (txErr error) {
		if unverifiedContactAlreadyExists {
			updateContactInput := contacts.UpdateContactInput{
				ID:                     contact.ID,
				SubscribedToNewsletter: opt.Bool(true),
				Country:                &country,
				FailedSignupAttempts:   opt.Int64(0),
				SignupCodeHash:         opt.String(codeHash),
			}
			txErr = service.contactsService.UpdateContactInternal(ctx, tx, &contact, updateContactInput)
			if txErr != nil {
				return txErr
			}
		} else {
			createContactInput := contacts.CreateContactInternalInput{
				Name:                   name,
				Email:                  email,
				Verified:               false,
				WebsiteID:              website.ID,
				Country:                country,
				SignupCodeHash:         codeHash,
				SubscribedToNewsletter: true,
			}
			contact, txErr = service.contactsService.CreateContactInternal(ctx, tx, createContactInput)
			if txErr != nil {
				return txErr
			}
		}

		job := queue.NewJobInput{
			Data: site.JobSendSubscribeEmail{
				Name:          contact.Name,
				Email:         contact.Email,
				Code:          code,
				ContactID:     contact.ID,
				WebsiteDomain: website.PrimaryDomain,
				WebsiteID:     website.ID,
			},
		}
		txErr = service.queue.Push(ctx, tx, job)
		if txErr != nil {
			errMessage := "site.Subscribe: Pushing job to queue"
			logger.Error(errMessage, slogx.Err(txErr))
			txErr = errs.Internal(errMessage, txErr)
			return txErr
		}

		return nil
	})
	if err != nil {
		return
	}

	ret.ContactID = contact.ID

	return
}
