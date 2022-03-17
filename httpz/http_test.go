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
	ctx := injector(context.Background())
	httpClient := httpz.Get(ctx)
	require.NotNil(t, httpClient)
}
