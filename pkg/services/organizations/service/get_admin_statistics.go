package service

import (
	"context"

	"golang.org/x/sync/errgroup"
	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/kernel"
	"markdown.ninja/pkg/services/organizations"
)

func (service *OrganizationsService) GetAdminStatistics(ctx context.Context, _ kernel.EmptyInput) (stats organizations.AdminStatistics, err error) {
	httpCtx := httpctx.FromCtx(ctx)

	accessToken := httpCtx.AccessToken
	if accessToken == nil || !accessToken.IsAdmin {
		return stats, kernel.ErrPermissionDenied
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	// for now we use a concurrency limit of 3 to balance between latency and database usage
	errGroup.SetLimit(3)

	errGroup.Go(func() error {
		var taskErr error
		stats.Organizations, taskErr = service.repo.GetOrganizationsCount(ctx, service.db)
		return taskErr
	})

	errGroup.Go(func() error {
		var taskErr error
		stats.PayingOrganizations, taskErr = service.repo.GetPayingOrganizationsCount(ctx, service.db)
		return taskErr
	})

	errGroup.Go(func() error {
		extraSlots, taskErr := service.repo.GetTotalExtraSlotsCount(ctx, service.db)
		if taskErr != nil {
			return taskErr
		}
		stats.MonthlyRevenue = extraSlots * kernel.PlanPro.Price
		return taskErr
	})

	err = errGroup.Wait()
	if err != nil {
		return stats, err
	}

	stats.MonthlyRevenue += (stats.PayingOrganizations * kernel.PlanPro.Price)

	return
}
