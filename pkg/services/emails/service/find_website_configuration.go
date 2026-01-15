package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/emails"
)

func (service *EmailsService) FindWebsiteConfiguration(ctx context.Context, db db.Queryer, websiteID guid.GUID) (configuration emails.WebsiteConfiguration, err error) {
	configuration, err = service.repo.FindWebsiteConfiguration(ctx, db, websiteID)
	return
}
