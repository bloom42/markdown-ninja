package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/uuid"
	"markdown.ninja/pkg/services/organizations"
)

func (service *OrganizationsService) ListOrganizationsForUser(ctx context.Context, db db.Queryer, userID uuid.UUID) (orgs []organizations.Organization, err error) {
	return service.repo.FindOrganizationsForUser(ctx, db, userID)
}
