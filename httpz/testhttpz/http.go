package testhttpz

import (
	"context"
	"net/http"
	"testing"

	"github.com/ibrt/golang-errors/errorz"
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

// Helper is a test helper for *http.Client.
type Helper struct {
	releaser func()
}

// BeforeSuite implements fixturez.BeforeSuite.
func (f *Helper) BeforeSuite(ctx context.Context, t *testing.T) context.Context {
	t.Helper()
	injector, releaser := httpz.Initializer(ctx)
	f.releaser = releaser
	return injector(ctx)
}

// AfterSuite implements fixturez.AfterSuite.
func (f *Helper) AfterSuite(_ context.Context, t *testing.T) {
	t.Helper()
	f.releaser()
	f.releaser = nil
}

// MockHelper is a test helper for *http.Client.
type MockHelper struct {
	releaser func()
}

// BeforeSuite implements fixturez.BeforeSuite.
func (f *MockHelper) BeforeSuite(ctx context.Context, t *testing.T) context.Context {
	t.Helper()
	injector, releaser := httpz.Initializer(ctx)
	f.releaser = releaser
	return injector(ctx)
}

// AfterSuite implements fixturez.AfterSuite.
func (f *MockHelper) AfterSuite(_ context.Context, t *testing.T) {
	t.Helper()
	f.releaser()
	f.releaser = nil
}

// BeforeTest implements fixturez.BeforeTest.
func (f *MockHelper) BeforeTest(ctx context.Context, t *testing.T) context.Context {
	t.Helper()
	gock.InterceptClient(httpz.Get(ctx))
	return ctx
}

// AfterTest implements fixturez.AfterTest.
func (f *MockHelper) AfterTest(ctx context.Context, t *testing.T) {
	t.Helper()

	if gock.HasUnmatchedRequest() {
		fixturez.AssertNoError(t, errorz.Errorf(
			"httpz.MockHelper: %v outgoing request(s) did not match any mock",
			errorz.A(len(gock.GetUnmatchedRequests())),
			errorz.M("unmatchedRequests", gock.GetUnmatchedRequests()),
			errorz.SkipPackage()))
	}

	if gock.IsPending() {
		fixturez.AssertNoError(t, errorz.Errorf(
			"httpz.MockHelper: %v mock(s) have not been matched by outgoing requests",
			errorz.A(len(gock.Pending())),
			errorz.M("pendingMocks", getPendingMocks(gock.Pending())),
			errorz.SkipPackage()))
	}

	gock.Flush()
	gock.CleanUnmatchedRequest()
	gock.RestoreClient(httpz.Get(ctx))
}

type pendingMock struct {
	Method     string
	URL        string
	PathParams map[string]string
	Header     http.Header
	Cookies    []*http.Cookie
	Body       string
	Counter    int
	Persisted  bool
}

func getPendingMocks(mocks []gock.Mock) []*pendingMock {
	pendingMocks := make([]*pendingMock, 0)

	for _, mock := range mocks {
		req := mock.Request()
		pendingMocks = append(pendingMocks, &pendingMock{
			Method: req.Method,
			URL: func() string {
				if req.URLStruct != nil {
					return req.URLStruct.String()
				}
				return ""
			}(),
			PathParams: func() map[string]string {
				if len(req.PathParams) > 0 {
					return req.PathParams
				}
				return nil
			}(),
			Header: func() http.Header {
				if len(req.Header) > 0 {
					return req.Header
				}
				return nil
			}(),
			Cookies: req.Cookies,
			Body: func() string {
				if len(req.BodyBuffer) > 1024 {
					return string(req.BodyBuffer)[:1024] + "..."
				}
				if len(req.BodyBuffer) > 0 {
					return string(req.BodyBuffer)
				}
				return ""
			}(),
			Counter:   req.Counter,
			Persisted: req.Persisted,
		})
	}

	return pendingMocks
}
