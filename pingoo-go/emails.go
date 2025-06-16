package pingoo

import (
	"context"
	"fmt"
	"net/http"
)

func (client *Client) CheckEmailAddress(ctx context.Context, email string) (res EmailInfo, err error) {
	err = client.request(ctx, requestParams{
		Method: http.MethodGet,
		Route:  fmt.Sprintf("/email/%s", email),
	}, &res)
	return
}

func (client *Client) LookupEmails(ctx context.Context, input LookupEmailsInput) (res []EmailInfo, err error) {
	res = make([]EmailInfo, 0, len(input.Emails))
	err = client.request(ctx, requestParams{
		Method:  http.MethodPost,
		Route:   "/lookup/emails",
		Payload: input,
	}, &res)
	return
}
