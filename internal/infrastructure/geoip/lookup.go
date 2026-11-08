package geoip

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type ipapiResponse struct {
	Country string `json:"country_name"`
	City    string `json:"city"`
}

type Resolver struct {
	client *http.Client
}

func New() *Resolver {
	return &Resolver{
		client: &http.Client{Timeout: 3 * time.Second},
	}
}

func (r *Resolver) Resolve(ctx context.Context, ipAddress string) (string, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "", fmt.Errorf("invalid client IP %q", ipAddress)
	}

	if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() {
		return "", nil
	}

	endpoint := "https://ipapi.co/" + url.PathEscape(ip.String()) + "/json/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("create GeoIP request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request GeoIP service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GeoIP service returned status %s", resp.Status)
	}

	var locationData ipapiResponse

	if err := json.NewDecoder(resp.Body).Decode(&locationData); err != nil {
		return "", fmt.Errorf("decode GeoIP response: %w", err)
	}

	switch {
	case locationData.Country != "" && locationData.City != "":
		return locationData.Country + ", " + locationData.City, nil
	case locationData.Country != "":
		return locationData.Country, nil
	default:
		return locationData.City, nil
	}
}
