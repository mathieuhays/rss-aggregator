package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetAuthorization(t *testing.T) {
	t.Run("invalid header format", func(t *testing.T) {
		header := http.Header{}
		header.Set("Authorization", "Bearer token onetoomanystring")

		_, err := GetAuthorization(header)
		if !errors.Is(err, errInvalid) {
			t.Errorf("expected 'invalid' error, got this instead: %v", err)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		header := http.Header{}
		header.Set("Authorization", "Test sometoken")

		_, err := GetAuthorization(header)
		if !errors.Is(err, errBadType) {
			t.Errorf("expected bad type error, got this instead: %v", err)
		}
	})

	t.Run("return valid token", func(t *testing.T) {
		header := http.Header{}
		tokenValue := "someKey"
		header.Set("Authorization", "Bearer "+tokenValue)

		auth, err := GetAuthorization(header)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if auth.Name != TypeBearer {
			t.Errorf("type mismatch. expected %s, got %s", TypeBearer, auth.Name)
		}

		if auth.Value != tokenValue {
			t.Errorf("value mismatch. expected %s, got %s", tokenValue, auth.Value)
		}
	})
}
