# golang-inject-http
[![Go Reference](https://pkg.go.dev/badge/github.com/ibrt/golang-inject-http.svg)](https://pkg.go.dev/github.com/ibrt/golang-inject-http)
![CI](https://github.com/ibrt/golang-inject-http/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/ibrt/golang-inject-http/branch/main/graph/badge.svg?token=BQVP881F9Z)](https://codecov.io/gh/ibrt/golang-inject-http)

HTTP client module for the [golang-inject](https://github.com/ibrt/golang-inject) framework.

### Basic Usage

This module injects a `*http.Client` into Go context. It provides both a real and a mock implementation for use in 
tests. Beyond being useful by itself, it is also a minimal example of how to tie together modules using the
[golang-inject](https://github.com/ibrt/golang-inject) framework, and how to easily test implementations using the
[golang-fixtures](https://github.com/ibrt/golang-fixtures) test suites.

```go
// main.go

package main

import (
    "io/ioutil"
    "net/http"
    
    "github.com/ibrt/golang-inject/injectz"
    "github.com/ibrt/golang-inject-http/httpz"
)

func main() {
    injector, releaser := injectz.Initialize(httpz.Initializer)
    defer releaser()
    
    middleware := injectz.NewMiddleware(injector)
    mux := http.NewServeMux()
    mux.Handle("/", middleware(http.HandlerFunc(handler)))
    _ = http.ListenAndServe(":3000", mux)
}

func handler(w http.ResponseWriter, r *http.Request) {
    httpResp, err := httpz.Get(r.Context()).Get("http://worldtimeapi.org/api/timezone/America/New_York")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    defer func() {
        _ = httpResp.Body.Close()
    }()
    
    buf, err := ioutil.ReadAll(httpResp.Body)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", httpResp.Header.Get("Content-Type"))
    w.WriteHeader(httpResp.StatusCode)
    _, _ = w.Write(buf)
}
```

```go
// main_test.go

package main

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/ibrt/golang-fixtures/fixturez"
    "github.com/ibrt/golang-inject-http/httpz/testhttpz"
    "github.com/stretchr/testify/require"
    "gopkg.in/h2non/gock.v1"
)

var (
    _ fixturez.Suite = &Suite{}
)

type Suite struct {
    *fixturez.DefaultConfigMixin
    HTTPZ *testhttpz.MockHelper
}

func TestSuite(t *testing.T) {
    fixturez.RunSuite(t, &Suite{})
}

func (s *Suite) TestHandler(ctx context.Context, t *testing.T) {
    const expectedResp = "{\"datetime\":\"2022-03-17T10:29:12.028398-04:00\"}\n"
    
    gock.New("http://worldtimeapi.org").
        Get("/api/timezone/America/New_York").
        Reply(200).
        JSON(json.RawMessage(expectedResp))
    
    w := httptest.NewRecorder()
    r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
    
    handler(w, r)
    
    require.Equal(t, http.StatusOK, w.Code)
    require.Equal(t, "application/json", w.Header().Get("Content-Type"))
    require.Equal(t, expectedResp, w.Body.String())
}
```

### Developers

Contributions are welcome, please check in on proposed implementation before sending a PR. You can validate your changes
using the `./test.sh` script.