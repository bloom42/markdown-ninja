package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/websites"
)

func (service *WebsitesService) FindRedirects(ctx context.Context, db db.Queryer, websiteID guid.GUID) (redirects []websites.Redirect, err error) {
	redirects, err = service.repo.FindRedirectsForWebsite(ctx, service.db, websiteID)
	return
}
