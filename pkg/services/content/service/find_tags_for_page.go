package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/content"
)

func (service *ContentService) FindTagsForPage(ctx context.Context, db db.Queryer, pageID guid.GUID) (tags []content.Tag, err error) {
	tags, err = service.repo.FindTagsForPage(ctx, service.db, pageID)
	if err != nil {
		return
	}

	return
}
