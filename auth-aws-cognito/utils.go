package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type Error struct {
	StatusCode int `json:"statusCode"`
	Message    any `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("API error: %d", e.StatusCode)
}

func Make(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var apiErr Error
			if errors.As(err, &apiErr) {
				err := WriteJson(w, apiErr.StatusCode, apiErr)
				if err != nil {
					return
				}
			}
		}
	}
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func StructToString(value interface{}) (string, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(value)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func StringToStruct(gobEncodedValue string, targetType interface{}) (interface{}, error) {
	typ := reflect.TypeOf(targetType)
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("targetType must be a struct type")
	}

	value := reflect.New(typ).Interface()

	reader := strings.NewReader(gobEncodedValue)
	err := gob.NewDecoder(reader).Decode(value)
	if err != nil {
		return nil, err
	}

	return reflect.ValueOf(value).Elem().Interface(), nil
}

func GetAuthProviderError(authErr error) string {
	parts := strings.Split(authErr.Error(), ",")

	lastPart := parts[len(parts)-1]
	finalParts := strings.Split(lastPart, ":")

	return strings.TrimSpace(finalParts[len(finalParts)-1])
}
