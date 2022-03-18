package testhttpz

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"

	"github.com/ibrt/golang-inject-http/httpz"
)

func TestHelpers(t *testing.T) {
	fixturez.RunSuite(t, &Suite{})
	fixturez.RunSuite(t, &MockSuite{})
}

type Suite struct {
	*fixturez.DefaultConfigMixin
	HTTPZ *Helper
}

func (s *Suite) TestHelper(ctx context.Context, t *testing.T) {
	httpClient := httpz.Get(ctx)
	require.NotNil(t, httpClient)
}

type MockSuite struct {
	*fixturez.DefaultConfigMixin
	HTTPZ *MockHelper
}

func (s *MockSuite) TestMockHelper(ctx context.Context, t *testing.T) {
	type Response struct {
		Value int `json:"value"`
	}

	respBody := &Response{
		Value: 10,
	}

	gock.New("https://mock-server.xyz").
		Get("/path").
		Reply(200).
		SetHeader("X-Custom-Header", "value").
		JSON(respBody)

	resp, err := httpz.Get(ctx).Get("https://mock-server.xyz/path")
	fixturez.RequireNoError(t, err)
	defer errorz.IgnoreClose(resp.Body)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "value", resp.Header.Get("X-Custom-Header"))

	respBodyBuf, err := ioutil.ReadAll(resp.Body)
	fixturez.RequireNoError(t, err)
	actualRespBody := &Response{}
	fixturez.RequireNoError(t, json.Unmarshal(respBodyBuf, actualRespBody))
	require.Equal(t, respBody, actualRespBody)

	_, err = httpz.Get(ctx).Get("https://mock-server.xyz/path")
	require.EqualError(t, err, `Get "https://mock-server.xyz/path": gock: cannot match any request`)

	_, err = httpz.Get(ctx).Get("https://another-mock-server.xyz")
	require.EqualError(t, err, `Get "https://another-mock-server.xyz": gock: cannot match any request`)
}

func TestGetPendingMocks(t *testing.T) {
	require.Equal(t,
		[]*pendingMock{
			{
				Counter: 1,
			},
			{
				Method:     "POST",
				PathParams: map[string]string{"k1": "v1"},
				Header:     http.Header{"K2": []string{"v2"}},
				Cookies:    []*http.Cookie{{Name: "test"}},
				Body:       "test",
				Counter:    2,
				Persisted:  true,
			},
			{
				Body:    strings.Repeat("a", 1024) + "...",
				Counter: 1,
			},
		},
		getPendingMocks([]gock.Mock{
			gock.NewMock(
				func() *gock.Request {
					req := gock.NewRequest()
					req.URLStruct = nil
					return req
				}(),
				gock.NewResponse()),
			gock.NewMock(
				func() *gock.Request {
					req := gock.NewRequest()
					req.Method = "POST"
					req.PathParams = map[string]string{"k1": "v1"}
					req.Header.Set("k2", "v2")
					req.Cookies = []*http.Cookie{{Name: "test"}}
					req.BodyBuffer = []byte("test")
					req.Counter = 2
					req.Persisted = true
					return req
				}(),
				gock.NewResponse()),
			gock.NewMock(
				func() *gock.Request {
					req := gock.NewRequest()
					req.BodyBuffer = []byte(strings.Repeat("a", 2048))
					return req
				}(),
				gock.NewResponse()),
		}))
}
