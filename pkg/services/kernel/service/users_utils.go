package service

import (
	"context"

	"github.com/skerkour/stdx-go/uuid"
	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/kernel"
)

func (service *KernelService) CurrentUserID(ctx context.Context) (userID uuid.UUID, err error) {
	httpCtx := httpctx.FromCtx(ctx)

	if httpCtx == nil || httpCtx.AccessToken == nil {
		err = kernel.ErrAuthenticationRequired
		return
	}

	return httpCtx.AccessToken.UserID, nil
}
