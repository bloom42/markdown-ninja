package pingoo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bloom42/stdx-go/retry"
)

type lookupHostInput struct {
	IpAddress string `json:"ip_address"`
}

type lookupHostOutput struct {
	Hostname string `json:"hostname"`
}

func (client *Client) resolveHostForIp(ctx context.Context, input lookupHostInput) (lookupHostOutput, error) {
	var hosts []string
	err := retry.Do(func() (retryErr error) {
		hosts, retryErr = client.dnsResolver.LookupAddr(ctx, input.IpAddress)
		if retryErr != nil {
			return retryErr
		}

		return nil
	}, retry.Context(ctx), retry.Attempts(4), retry.Delay(50*time.Millisecond))
	if err != nil {
		return lookupHostOutput{}, fmt.Errorf("waf: error resolving hosts for IP address (%s): %w", input.IpAddress, err)
	}

	cleanedUpHosts := make([]string, 0, len(hosts))
	for _, host := range hosts {
		host = strings.ToValidUTF8(strings.TrimSuffix(strings.TrimSpace(host), "."), "")
		if host != "" {
			cleanedUpHosts = append(cleanedUpHosts, host)
		}
	}
	hosts = cleanedUpHosts

	if len(hosts) > 0 {
		return lookupHostOutput{Hostname: hosts[0]}, nil
	}

	return lookupHostOutput{Hostname: ""}, nil
}
