package storehost

import "testing"

func TestSlugFromHost(t *testing.T) {
	cases := []struct {
		host string
		base string
		want string
	}{
		{"north.localhost:5173", "localhost", "north"},
		{"localhost:5173", "localhost", ""},
		{"www.example.com", "example.com", ""},
		{"north.example.com", "example.com", "north"},
		{"a.b.example.com", "example.com", ""},
		{"example.com", "example.com", ""},
	}
	for _, tc := range cases {
		if got := SlugFromHost(tc.host, tc.base); got != tc.want {
			t.Fatalf("SlugFromHost(%q, %q) = %q, want %q", tc.host, tc.base, got, tc.want)
		}
	}
}
