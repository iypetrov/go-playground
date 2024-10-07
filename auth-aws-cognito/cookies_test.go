package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteAndReadCookie(t *testing.T) {
	secretKey := []byte(os.Getenv("APP_SECRET"))
	expected := UserCookie{
		Email:        "user@example.com",
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_123",
	}

	rr := httptest.NewRecorder()

	err := WriteCookie(rr, "APP_COOKIE", expected, secretKey)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := rr.Result().Cookies()[0]
	req.AddCookie(cookie)

	// actual, err := ReadCookie[RegisterResponse](req, "APP_COOKIE", secretKey)
	// require.NoError(t, err)

	// require.Equal(t, expected.Email, actual.Email)
	// require.Equal(t, expected.AccessToken, actual.AccessToken)
	// require.Equal(t, expected.RefreshToken, actual.RefreshToken)
}
