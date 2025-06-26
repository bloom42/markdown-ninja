package service

import (
	"context"

	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/kernel"
	"markdown.ninja/pkg/services/websites"
)

func (service *WebsitesService) GetAdminStatistics(ctx context.Context, _ kernel.EmptyInput) (stats websites.AdminStatistics, err error) {
	httpCtx := httpctx.FromCtx(ctx)

	accessToken := httpCtx.AccessToken
	if accessToken == nil || !accessToken.IsAdmin {
		return stats, kernel.ErrPermissionDenied
	}

	stats.Websites, err = service.repo.GetWebsitesCount(ctx, service.db)
	if err != nil {
		return stats, err
	}

	return stats, nil
}
