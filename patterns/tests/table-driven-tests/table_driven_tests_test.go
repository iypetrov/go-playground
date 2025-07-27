package tabledriventests

import "testing"

func TestAddition(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "foo",
			a:        1,
			b:        1,
			expected: 2,
		},
		{
			name:     "bar",
			a:        2,
			b:        3,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
