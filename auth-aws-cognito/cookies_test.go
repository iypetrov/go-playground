package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteAndReadCookie(t *testing.T) {
    secretKey := []byte("13d6b4dff8f84a10851021ec8608f814")

    inputCookie := UserCookie{
        Email:        "test@example.com",
        AccessToken:  "access-token",
        RefreshToken: "refresh-token",
    }

    w := httptest.NewRecorder()
    err := WriteCookie(w, "user_cookie", inputCookie, secretKey)
    assert.NoError(t, err) 

    cookie := w.Result().Cookies()
    assert.Len(t, cookie, 1)
    assert.Equal(t, "APP_COOKIE", cookie[0].Name)

    r := httptest.NewRequest(http.MethodGet, "/", nil)
    r.AddCookie(cookie[0]) 

    outputCookie, err := ReadCookie(r, "APP_COOKIE", secretKey)
    assert.NoError(t, err) 

    assert.Equal(t, inputCookie, outputCookie)
}
