package testhttpz

import (
	"context"
	"testing"

	"github.com/ibrt/golang-fixtures/fixturez"
	"gopkg.in/h2non/gock.v1"

	"github.com/ibrt/golang-inject-http/httpz"
)

var (
	_ fixturez.BeforeSuite = &Helper{}
	_ fixturez.AfterSuite  = &Helper{}
	_ fixturez.BeforeSuite = &MockHelper{}
	_ fixturez.AfterSuite  = &MockHelper{}
	_ fixturez.BeforeTest  = &MockHelper{}
	_ fixturez.AfterTest   = &MockHelper{}
)

// Helper provides a test helper for httpz using a real http.Transport.
type Helper struct {
	releaser func()
}

// BeforeSuite implements fixturez.BeforeSuite.
func (f *Helper) BeforeSuite(ctx context.Context, _ *testing.T) context.Context {
	injector, releaser := httpz.Initializer(ctx)
	f.releaser = releaser
	return injector(ctx)
}

// AfterSuite implements fixturez.AfterSuite.
func (f *Helper) AfterSuite(_ context.Context, _ *testing.T) {
	f.releaser()
	f.releaser = nil
}

// MockHelper provides a test helper for httpz using a mock http.Transport.
type MockHelper struct {
	releaser func()
}

// BeforeSuite implements fixturez.BeforeSuite.
func (f *MockHelper) BeforeSuite(ctx context.Context, _ *testing.T) context.Context {
	injector, releaser := httpz.Initializer(ctx)
	f.releaser = releaser
	return injector(ctx)
}

// AfterSuite implements fixturez.AfterSuite.
func (f *MockHelper) AfterSuite(_ context.Context, _ *testing.T) {
	f.releaser()
	f.releaser = nil
}

// BeforeTest implements fixturez.BeforeTest.
func (f *MockHelper) BeforeTest(ctx context.Context, _ *testing.T) context.Context {
	gock.InterceptClient(httpz.Get(ctx))
	return ctx
}

// AfterTest implements fixturez.AfterTest.
func (f *MockHelper) AfterTest(ctx context.Context, _ *testing.T) {
	gock.Off()
	gock.RestoreClient(httpz.Get(ctx))
}
