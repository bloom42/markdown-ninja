package service

import (
	"context"
	"math"

	"github.com/skerkour/stdx-go/db"
	"markdown.ninja/pkg/services/websites"
)

func (service *WebsitesService) FindAllWebsites(ctx context.Context, db db.Queryer) (websites []websites.Website, err error) {
	return service.repo.FindWebsites(ctx, service.db, math.MaxInt64)
}
