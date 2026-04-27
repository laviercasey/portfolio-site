package service

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

var blockedMetadataHosts = map[string]struct{}{
	"169.254.169.254":          {},
	"metadata.google.internal": {},
	"metadata.goog":            {},
	"metadata":                 {},
	"169.254.170.2":            {},
}

func validateUmamiURL(raw, env string) (*url.URL, error) {
	if raw == "" {
		return nil, fmt.Errorf("UMAMI_API_URL is empty")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse UMAMI_API_URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid scheme %q", u.Scheme)
	}
	host := u.Hostname()
	if host == "" {
		return nil, fmt.Errorf("missing host")
	}
	lower := strings.ToLower(host)
	if _, blocked := blockedMetadataHosts[lower]; blocked {
		return nil, fmt.Errorf("host %q is a cloud metadata endpoint and is not allowed", host)
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLinkLocalUnicast() {
		return nil, fmt.Errorf("link-local host %q is not allowed", host)
	}
	if env != "dev" {
		loopback := map[string]bool{
			"localhost": true,
			"127.0.0.1": true,
			"::1":       true,
			"0.0.0.0":   true,
		}
		if loopback[lower] {
			return nil, fmt.Errorf("loopback host %q is not allowed outside dev", host)
		}
		if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
			return nil, fmt.Errorf("loopback host %q is not allowed outside dev", host)
		}
	}
	return u, nil
}
