package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
)

func (service *WebsitesService) GetDomainsCountForWebsite(ctx context.Context, db db.Queryer, websiteID guid.GUID) (count int64, err error) {
	return service.repo.GetDomainsCountForWebsite(ctx, db, websiteID)
}
