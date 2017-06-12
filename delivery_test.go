package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type fmturltestpair struct {
	data   postback
	result string
}

type httptestpair struct {
  url string
  requestType string
  result responseData
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

var httpSendTests = []httptestpair {
  {url: "https://httpbin.org/get?evil=100+dollars", requestType: "GET", result: responseData{responseCode: "200", responseTime: "???", responseBody: "???"}},

}

func TestFormatUrl(t *testing.T) {
	for _, pair := range urlTests {
		r := formatUrl(pair.data)
		assert.Equal(t, r.Url, pair.result, "formatUrl() didn't return the expected formatted url.")
	}
}

func TestSendRequest(t *testing.T) {
	for _, obj := range httpSendTests {
		r, _ := sendRequest(obj.url, obj.requestType)
		assert.Equal(t, obj.result.responseCode, r.responseCode, "sendRequest() didn't return expected response code.")
	}
}
