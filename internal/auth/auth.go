package auth

import (
	"errors"
	"net/http"
	"slices"
	"strings"
)

const (
	TypeBearer = "Bearer"
	TypeApiKey = "ApiKey"
)

var (
	errInvalid = errors.New("invalid authorization header")
	errBadType = errors.New("authorization type not recognised")
)

var validTypes = []string{
	TypeBearer,
	TypeApiKey,
}

type Authorization struct {
	Name  string
	Value string
}

func GetAuthorization(header http.Header) (Authorization, error) {
	parts := strings.Split(header.Get("Authorization"), " ")
	if len(parts) != 2 {
		return Authorization{}, errInvalid
	}

	if !slices.Contains(validTypes, parts[0]) {
		return Authorization{}, errBadType
	}

	return Authorization{
		Name:  parts[0],
		Value: parts[1],
	}, nil
}
