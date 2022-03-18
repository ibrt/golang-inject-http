package httpz_test

import (
	"context"
	"testing"

	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"

	"github.com/ibrt/golang-inject-http/httpz"
)

func TestModule(t *testing.T) {
	injector, releaser := httpz.Initializer(context.Background())
	defer releaser()
	ctx := injector(context.Background())
	httpClient := httpz.Get(ctx)
	require.NotNil(t, httpClient)
	require.Nil(t, httpz.MaybeGet(context.Background()))
	fixturez.RequirePanicsWith(t, "httpz: not initialized", func() {
		httpz.Get(context.Background())
	})
}
