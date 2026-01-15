package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
	"markdown.ninja/pkg/services/content"
)

func (service *ContentService) FindSnippets(ctx context.Context, db db.Queryer, websiteID guid.GUID) (snippets []content.Snippet, err error) {
	snippets, err = service.repo.FindSnippetsForWebsite(ctx, db, websiteID)
	return
}
