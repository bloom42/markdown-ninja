package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/content"
)

func (service *ContentService) FindLastPublishedPageOrPost(ctx context.Context, db db.Queryer, websiteID guid.GUID) (page content.Page, err error) {
	page, err = service.repo.FindLastPublishedPageOrPostForWebsite(ctx, db, websiteID)
	return
}
