package pingoo

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"
)

func (client *Client) LookupEmail(ctx context.Context, email string) (res EmailInfo, err error) {
	err = client.request(ctx, requestParams{
		Method: http.MethodGet,
		Route:  fmt.Sprintf("/lookup/email/%s", email),
	}, &res)
	return
}

func (client *Client) LookupIp(ctx context.Context, ipAddress netip.Addr) (res IpInfo, err error) {
	err = client.request(ctx, requestParams{
		Method: http.MethodGet,
		Route:  fmt.Sprintf("/lookup/ip/%s", ipAddress),
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
