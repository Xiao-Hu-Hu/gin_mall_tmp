package util

import (
	"net/http"
	"strings"
)

// ExtractTokenFromRequest extracts token from token or Authorization header.
func ExtractTokenFromRequest(r *http.Request) string {
	token := strings.TrimSpace(r.Header.Get("token"))
	if token != "" {
		return token
	}

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}

	return authHeader
}
