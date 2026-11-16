package storehost

import (
	"net"
	"strings"

	"shop24/internal/domain/input"
)

// SlugFromHost returns the store slug for a host such as north.localhost or north.example.com.
// The apex domain and www host are not storefronts.
func SlugFromHost(hostport, baseDomain string) string {
	host := hostport
	if parsed, _, err := net.SplitHostPort(hostport); err == nil {
		host = parsed
	}
	host = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
	base := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(baseDomain)), ".")
	if host == "" || base == "" || host == base || host == "www."+base {
		return ""
	}
	suffix := "." + base
	if !strings.HasSuffix(host, suffix) {
		return ""
	}
	slug := strings.TrimSuffix(host, suffix)
	if !input.ValidSlug(slug) {
		return ""
	}
	return slug
}
