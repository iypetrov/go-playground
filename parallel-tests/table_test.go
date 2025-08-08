package main

import (
	"testing"
	"time"
)

func TestApi_with_test_table(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name string
		API  string
	}{
		{Name: "1", API: "/api/1"},
		{Name: "2", API: "/api/2"},
		{Name: "3", API: "/api/3"},
		{Name: "4", API: "/api/4"},
		{Name: "5", API: "/api/5"},
		{Name: "6", API: "/api/6"},
		{Name: "7", API: "/api/7"},
		{Name: "8", API: "/api/8"},
		{Name: "9", API: "/api/9"},
		{Name: "10", API: "/api/9"},
		{Name: "11", API: "/api/1"},
		{Name: "12", API: "/api/2"},
		{Name: "13", API: "/api/3"},
		{Name: "14", API: "/api/4"},
		{Name: "15", API: "/api/5"},
		{Name: "16", API: "/api/6"},
		{Name: "17", API: "/api/7"},
		{Name: "18", API: "/api/8"},
		{Name: "19", API: "/api/9"},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			simulateSlowCall(1 * time.Second)
		})
	}
}
