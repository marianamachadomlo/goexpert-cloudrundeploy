package cep

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		name    string
		zipcode string
		want    bool
	}{
		{name: "valid zipcode", zipcode: "01310100", want: true},
		{name: "too short", zipcode: "1234567", want: false},
		{name: "too long", zipcode: "123456789", want: false},
		{name: "contains hyphen", zipcode: "01310-100", want: false},
		{name: "contains letters", zipcode: "01310abc", want: false},
		{name: "empty", zipcode: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.zipcode); got != tt.want {
				t.Fatalf("IsValid(%q) = %v, want %v", tt.zipcode, got, tt.want)
			}
		})
	}
}
