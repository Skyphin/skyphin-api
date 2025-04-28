package utils

import (
	"net/url"
	"strings"
)

type URLService struct {
	allowedDomains map[string]bool
}

func NewURLService(allowedDomains []string) *URLService {
	service := &URLService{
		allowedDomains: make(map[string]bool),
	}
	for _, domain := range allowedDomains {
		service.allowedDomains[domain] = true
	}
	return service
}

func (s *URLService) Normalize(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Force HTTPS and remove www
	u.Scheme = "https"
	u.Host = strings.TrimPrefix(u.Host, "www.")

	// Validate domain if we have restrictions
	if len(s.allowedDomains) > 0 {
		domainParts := strings.Split(u.Host, ".")
		if len(domainParts) < 2 {
			return "", err
		}
		domain := strings.Join(domainParts[len(domainParts)-2:], ".")
		if !s.allowedDomains[domain] {
			return "", err
		}
	}

	// Normalize path
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawQuery = ""
	u.Fragment = ""

	return u.String(), nil
}
