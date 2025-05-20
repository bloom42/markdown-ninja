package service

import (
	"context"
	"net/http"
	"time"

	"github.com/bloom42/stdx-go/db"
	"github.com/bloom42/stdx-go/log/slogx"
	"markdown.ninja/pkg/errs"
	"markdown.ninja/pkg/services/contacts"
)

func (service *ContactsService) MarkSessionAsVerified(ctx context.Context, db db.Queryer, session *contacts.Session) (sessionCookie http.Cookie, err error) {
	logger := slogx.FromCtx(ctx)

	if session == nil {
		errMessage := "contacts.MarkSessionAsVerified: session is null"
		logger.Error(errMessage)
		err = errs.Internal(errMessage, nil)
		return
	}

	newSessionToken, err := service.generateSessionToken(session.ContactID, session.ID)
	if err != nil {
		return
	}

	session.UpdatedAt = time.Now().UTC()
	session.CodeHash = ""
	session.Verified = true
	session.SecretHash = newSessionToken.hash[:]

	err = service.repo.UpdateSession(ctx, db, *session)
	if err != nil {
		return
	}

	sessionCookie = service.generateSessionCookie(newSessionToken.token)
	return
}
