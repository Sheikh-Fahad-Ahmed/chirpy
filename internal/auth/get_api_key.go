package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error){
	authHeaders := headers.Get("Authorization")
	if authHeaders == "" {
		return "", errors.New("header does not exist")
	}
	parts := strings.Fields(authHeaders)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "apikey" {
		return "", errors.New("invalid authorization header")
	}
	return parts[1], nil
}