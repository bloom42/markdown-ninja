package service

import (
	"context"

	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/kernel"
)

func (service *SiteService) Logout(ctx context.Context, input kernel.EmptyInput) (err error) {
	httpCtx := httpctx.FromCtx(ctx)

	contact := service.contactsService.CurrentContact(ctx)
	if contact == nil {
		err = kernel.ErrAuthenticationRequired
		return
	}

	session := service.contactSession(ctx)

	err = service.contactsService.DeleteSession(ctx, service.db, session.ID)
	if err != nil {
		return
	}

	logoutCookie := service.contactsService.GenerateLogoutCookie()
	httpCtx.Response.Cookies = append(httpCtx.Response.Cookies, logoutCookie)
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Clear-Site-Data
	httpCtx.Response.Headers.Set("Clear-Site-Data", "*")

	return
}
