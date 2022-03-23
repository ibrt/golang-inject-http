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

// Initializer is a *http.Client initializer.
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

// NewSingletonInjector always injects the given *http.Client.
func NewSingletonInjector(httpClient *http.Client) injectz.Injector {
	return injectz.NewSingletonInjector(httpClientContextKey, httpClient)
}

// Get extracts the *http.Client from context, returns the http.DefaultClient if not found.
func Get(ctx context.Context) *http.Client {
	if httpClient, ok := ctx.Value(httpClientContextKey).(*http.Client); ok {
		return httpClient
	}
	return http.DefaultClient
}
