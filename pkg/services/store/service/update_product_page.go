package service

import (
	"context"
	"strings"
	"time"

	"github.com/skerkour/stdx-go/crypto/blake3"
	"markdown.ninja/pkg/services/store"
)

func (service *StoreService) UpdateProductPage(ctx context.Context, input store.UpdateProductPageInput) (page store.ProductPage, err error) {
	actorID, err := service.kernel.CurrentUserID(ctx)
	if err != nil {
		return
	}

	page, err = service.repo.FindProductPageByID(ctx, service.db, input.ID)
	if err != nil {
		return
	}

	product, err := service.repo.FindProductByID(ctx, service.db, page.ProductID)
	if err != nil {
		return
	}

	err = service.websitesService.CheckUserIsStaff(ctx, service.db, actorID, product.WebsiteID)
	if err != nil {
		return
	}

	now := time.Now().UTC()
	page.UpdatedAt = now

	if input.Title != nil {
		page.Title = strings.TrimSpace(*input.Title)
		err = service.validateProductPageTitle(page.Title)
		if err != nil {
			return
		}
	}

	if input.BodyMarkdown != nil {
		page.BodyMarkdown = *input.BodyMarkdown
		err = service.validateProductPageContent(page.BodyMarkdown)
		if err != nil {
			return
		}

		page.Size = int64(len(page.BodyMarkdown))
		bodyHash := blake3.Sum256([]byte(page.BodyMarkdown))
		page.Hash = bodyHash[:]
	}

	err = service.repo.UpdateProductPage(ctx, service.db, page)
	if err != nil {
		return
	}

	return
}
