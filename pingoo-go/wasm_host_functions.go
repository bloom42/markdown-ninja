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

func (client *Client) resolveHostForIp(ctx context.Context, ip string) (string, error) {
	logger := client.getLogger(ctx)

	var hosts []string
	err := retry.Do(func() (retryErr error) {
		hosts, retryErr = client.dnsResolver.LookupAddr(ctx, ip)
		return retryErr
	}, retry.Context(ctx), retry.Attempts(5), retry.Delay(10*time.Millisecond))
	if err != nil {
		logger.Warn(fmt.Sprintf("pingoo: error resolving hostname for IP [%s]: %s", ip, err.Error()))
		return "", nil
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
		return hosts[0], nil
	}

	return "", nil
}
