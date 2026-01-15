package service

import (
	"context"
	"time"

	"github.com/skerkour/stdx-go/db"
	"markdown.ninja/pkg/services/contacts"
)

func (service *ContactsService) FailSessionLoginAttempt(ctx context.Context, db db.Queryer, session contacts.Session) (err error) {
	session.UpdatedAt = time.Now().UTC()
	session.FailedLoginAttempts += 1
	err = service.repo.UpdateSession(ctx, service.db, session)
	return
}
