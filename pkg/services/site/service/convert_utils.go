package service

import (
	"context"
	"html/template"
	"time"

	"markdown.ninja/pkg/services/contacts"
	"markdown.ninja/pkg/services/content"
	"markdown.ninja/pkg/services/site"
	"markdown.ninja/pkg/services/store"
	"markdown.ninja/pkg/services/websites"
)

func (service *SiteService) convertContact(input contacts.Contact) site.Contact {
	subscribedToNewsletter := false
	if input.SubscribedToNewsletterAt != nil {
		subscribedToNewsletter = true
	}

	return site.Contact{
		Name:                   input.Name,
		Email:                  input.Email,
		SubscribedToNewsletter: subscribedToNewsletter,
	}
}

func (service *SiteService) convertOrder(input store.Order) site.Order {
	return site.Order{
		ID:          input.ID,
		CreatedAt:   input.CreatedAt.Truncate(time.Second),
		TotalAmount: input.TotalAmount,
		Currency:    input.Currency,
		Status:      input.Status,
		InvoiceUrl:  input.StripeInvoiceUrl,
	}
}

func (service *SiteService) convertOrders(input []store.Order) []site.Order {
	ret := make([]site.Order, len(input))

	for i, item := range input {
		ret[i] = service.convertOrder(item)
	}
	return ret
}

func (service *SiteService) convertProduct(website websites.Website, input store.Product) (ret site.Product) {
	pages := service.convertProductPages(website, input.Content)

	ret = site.Product{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,

		Content: pages,
	}
	return ret
}

func (service *SiteService) convertProducts(website websites.Website, input []store.Product) (ret []site.Product) {
	ret = make([]site.Product, len(input))

	for i, item := range input {
		ret[i] = service.convertProduct(website, item)
	}

	return ret
}

func (service *SiteService) convertProductPage(website websites.Website, input store.ProductPage) (ret site.ProductPage) {
	ret = site.ProductPage{
		ID:       input.ID,
		Position: input.Position,
		Title:    input.Title,
		Body:     service.contentService.RenderMarkdown(website, input.BodyMarkdown, nil, false),
	}
	return ret
}

func (service *SiteService) convertProductPages(website websites.Website, input []store.ProductPage) (ret []site.ProductPage) {
	if input == nil {
		return ret
	}

	ret = make([]site.ProductPage, len(input))

	for i, item := range input {
		ret[i] = service.convertProductPage(website, item)
	}

	return ret
}

func (service *SiteService) convertWebsite(input websites.Website) site.Website {
	url := service.httpConfig.WebsitesBaseUrl.Scheme + "://" + input.PrimaryDomain + service.httpConfig.WebsitesPort

	return site.Website{
		Url:          template.URL(url),
		Name:         input.Name,
		Description:  input.Description,
		Navigation:   input.Navigation,
		Language:     input.Language,
		Ad:           input.Ad,
		Announcement: input.Announcement,
		Colors:       input.Colors,
		Logo:         input.Logo,
		PoweredBy:    input.PoweredBy,
		Theme:        input.Theme,
	}
}

func (service *SiteService) convertPage(_ context.Context, website websites.Website, input content.Page, tags []content.Tag, snippets []content.Snippet) (ret site.Page) {
	if tags == nil {
		tags = []content.Tag{}
	}

	bodyHtml := service.contentService.RenderMarkdown(website, input.BodyMarkdown, snippets, false)

	ret = site.Page{
		PageMetadata: service.convertPageToMetadata(website, input),
		Tags:         service.convertTags(tags),
		Body:         bodyHtml,
	}
	return ret
}

func (service *SiteService) convertTags(input []content.Tag) []site.Tag {
	ret := make([]site.Tag, len(input))

	for i, tag := range input {
		ret[i] = site.Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}
	}

	return ret
}

func (service *SiteService) convertPageToMetadata(website websites.Website, page content.Page) site.PageMetadata {
	url := service.httpConfig.WebsitesBaseUrl.Scheme + "://" + website.PrimaryDomain + service.httpConfig.WebsitesPort + page.Path

	return site.PageMetadata{
		Date:         page.Date.UTC().Truncate(time.Minute),
		ModifiedAt:   page.ModifiedAt().UTC().Truncate(time.Minute),
		Type:         page.Type,
		Title:        page.Title,
		Path:         page.Path,
		Description:  page.Description,
		Language:     page.Language,
		Url:          template.URL(url),
		BodyHash:     page.BodyHash,
		MetadataHash: page.MetadataHash,
	}
}

func (service *SiteService) convertPageMetadata(website websites.Website, page content.PageMetadata) site.PageMetadata {
	url := service.httpConfig.WebsitesBaseUrl.Scheme + "://" + website.PrimaryDomain + service.httpConfig.WebsitesPort + page.Path

	return site.PageMetadata{
		Date:         page.Date.UTC().Truncate(time.Minute),
		ModifiedAt:   page.ModifiedAt().UTC().Truncate(time.Minute),
		Type:         page.Type,
		Title:        page.Title,
		Path:         page.Path,
		Description:  page.Description,
		Language:     page.Language,
		Url:          template.URL(url),
		BodyHash:     page.BodyHash,
		MetadataHash: page.MetadataHash,
	}
}

// I know that metadata is already plural... :)
func (service *SiteService) convertPageMetadatas(website websites.Website, input []content.PageMetadata) []site.PageMetadata {
	ret := make([]site.PageMetadata, len(input))

	for i, elem := range input {
		ret[i] = service.convertPageMetadata(website, elem)
	}

	return ret
}

func generatePageBodyHtmlCacheKey(page content.Page) string {
	return page.BodyHash.String()
}
