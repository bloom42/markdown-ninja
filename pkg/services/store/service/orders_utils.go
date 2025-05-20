package service

import (
	"fmt"

	"github.com/bloom42/stdx-go/guid"
	"markdown.ninja/pkg/services/store"
)

func (service *StoreService) generateCompleteOrderUrl(domain string, orderID guid.GUID) string {
	hostname := domain + service.websitesPort
	return fmt.Sprintf("%s://%s/checkout/%s/complete",
		service.httpConfig.WebsitesBaseUrl.Scheme, hostname, orderID.String())
}

func (service *StoreService) generateCancelOrderUrl(domain string, orderID guid.GUID) string {
	hostname := domain + service.websitesPort
	return fmt.Sprintf("%s://%s/checkout/%s/cancel",
		service.httpConfig.WebsitesBaseUrl.Scheme, hostname, orderID.String())
}

func convertOrderToMetadata(order store.Order) store.OrderMetadata {
	return store.OrderMetadata{
		ID:          order.ID,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		Email:       order.Email,
		TotalAmount: order.TotalAmount,
		Currency:    order.Currency,
		Status:      order.Status,
		CompletedAt: order.CompletedAt,
		CanceledAt:  order.CanceledAt,
		ContactID:   order.ContactID,
	}
}
