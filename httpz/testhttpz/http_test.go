package testhttpz_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"

	"github.com/ibrt/golang-inject-http/httpz"
	"github.com/ibrt/golang-inject-http/httpz/testhttpz"
)

func TestHelpers(t *testing.T) {
	fixturez.RunSuite(t, &Suite{})
	fixturez.RunSuite(t, &MockSuite{})
}

type Suite struct {
	*fixturez.DefaultConfigMixin
	HTTPZ *testhttpz.Helper
}

func (s *Suite) TestHelper(ctx context.Context, t *testing.T) {
	httpClient := httpz.Get(ctx)
	require.NotNil(t, httpClient)
}

type MockSuite struct {
	*fixturez.DefaultConfigMixin
	HTTPZ *testhttpz.MockHelper
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
}
