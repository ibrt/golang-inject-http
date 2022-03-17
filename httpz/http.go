package httpz

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/ibrt/golang-inject/injectz"
)

type contextKey int

const (
	httpClientContextKey contextKey = iota
)

var (
	_ injectz.Initializer = Initializer
)

// Initializer is a HTTP client initializer.
func Initializer(_ context.Context) (injectz.Injector, injectz.Releaser) {
	return NewSingletonInjector(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}), injectz.NewNoopReleaser()
}

// NewSingletonInjector always injects the given HTTP client.
func NewSingletonInjector(httpClient *http.Client) injectz.Injector {
	return injectz.NewSingletonInjector(httpClientContextKey, httpClient)
}

// Get returns the HTTP client, panics if not found.
func Get(ctx context.Context) *http.Client {
	return ctx.Value(httpClientContextKey).(*http.Client)
}
