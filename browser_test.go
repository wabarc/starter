package starter

import (
	"os"
	"strings"
	"testing"
)

func TestFlags(t *testing.T) {
	b := &Browser{}

	tests := []struct {
		key  string
		val  string
		want bool
	}{
		{
			key:  "",
			val:  "",
			want: false,
		},
		{
			key:  "",
			val:  "invalid",
			want: false,
		},
		{
			key:  "http_proxy",
			val:  "http://127.0.0.1:1080",
			want: true,
		},
		{
			key:  "HTTP_PROXY",
			val:  "http://127.0.0.1:1080",
			want: true,
		},
		{
			key:  "HTTPS_PROXY",
			val:  "https://127.0.0.1:1080",
			want: true,
		},
		{
			key:  "HTTP_PROXY",
			val:  "socks5://127.0.0.1:1080",
			want: true,
		},
	}
	for _, test := range tests {
		os.Clearenv()
		os.Setenv(test.key, test.val)
		flags := b.flags()
		var hasProxy bool
		for _, flag := range flags {
			hasProxy = strings.HasPrefix(flag, "--proxy-server")
		}
		if hasProxy != test.want {
			t.Fatalf("unexpected set proxy flag got %t instead of %t", hasProxy, test.want)
		}
	}
}
