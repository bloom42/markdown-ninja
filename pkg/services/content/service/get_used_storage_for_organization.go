package service

import (
	"context"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/guid"
)

func (service *ContentService) GetUsedStorageForOrganization(ctx context.Context, db db.Queryer, organizationID guid.GUID) (storage int64, err error) {
	return service.repo.GetUsedAssetsStorageForOrganization(ctx, db, organizationID)
}
