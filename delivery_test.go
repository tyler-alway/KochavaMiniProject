package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fmturltestpair struct {
	data   postback
	result string
}

var urlTests = []fmturltestpair{
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={$money}", Data: map[string]string{"$money": "100 dollars"}}, result: "https://httpbin.org/get?evil=100+dollars"},
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={$money}"}, result: "https://httpbin.org/get?evil="},
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil="}, result: "https://httpbin.org/get?evil="},
	{data: postback{Method: "GET", Url: "", Data: map[string]string{"$money": "100 dollars"}}, result: ""},
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil="}, result: "https://httpbin.org/get?evil="},
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={$money}&other={other}", Data: map[string]string{"$money": "100 dollars", "other": "something else"}}, result: "https://httpbin.org/get?evil=100+dollars&other=something+else"},
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={$money}&other={other}"}, result: "https://httpbin.org/get?evil=&other="},
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={$money}&other={other}", Data: map[string]string{"other": "something else"}}, result: "https://httpbin.org/get?evil=&other=something+else"},
	{data: postback{Method: "GET", Url: "{.*}", Data: map[string]string{".*": "100 dollars", "other": "something else"}}, result: "100+dollars"},
	{data: postback{Method: "GET", Url: "{$money*}{other?}", Data: map[string]string{"$money*": "100 dollars", "other?": "something else"}}, result: "100+dollarssomething+else"},
}

func TestFormatUrl(t *testing.T) {
	for _, pair := range urlTests {
		r := formatUrl(pair.data)
		assert.Equal(t, pair.result, r.Url, "formatUrl() didn't return the expected formatted url.")
	}
}

func TestSendRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "body string")
	}))
	defer testServer.Close()

	testURl := testServer.URL

	resp, _ := sendRequest(testURl, "GET")

	assert.Equal(t, "200", resp.responseCode, "sendRequest() didn't return the expected status code.")
	assert.Equal(t, "body string\n", resp.responseBody, "sendRequest() didn't return the expected response body.")
}
