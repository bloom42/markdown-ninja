package service

import (
	"context"
	"fmt"

	"github.com/bloom42/stdx-go/guid"
	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/contacts"
)

func (service *SiteService) contactSession(ctx context.Context) (session *contacts.Session) {
	httpCtx := httpctx.FromCtx(ctx)

	return httpCtx.ContactSession
}

func (service *SiteService) generateLoginLink(domain string, sessionID guid.GUID, code string) (link string) {
	hostname := domain + service.httpConfig.WebsitesPort
	link = fmt.Sprintf("%s://%s/login?session=%s&code=%s",
		service.httpConfig.WebsitesBaseUrl.Scheme, hostname, sessionID.String(), code)
	return
}

func (service *SiteService) generateSubscribeLink(domain string, contactID guid.GUID, code string) (link string) {
	hostname := domain + service.httpConfig.WebsitesPort
	link = fmt.Sprintf("%s://%s/subscribe?contact=%s&code=%s",
		service.httpConfig.WebsitesBaseUrl.Scheme, hostname, contactID.String(), code)
	return
}
