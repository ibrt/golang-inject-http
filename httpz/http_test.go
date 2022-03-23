package httpz_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ibrt/golang-inject-http/httpz"
)

func TestModule(t *testing.T) {
	injector, releaser := httpz.Initializer(context.Background())
	defer releaser()

	httpClient := httpz.Get(context.Background())
	require.NotNil(t, httpClient)

	httpClient = httpz.Get(injector(context.Background()))
	require.NotNil(t, httpClient)

	httpClient = httpz.Get(httpz.NewSingletonInjector(httpClient)(context.Background()))
	require.NotNil(t, httpClient)
}
