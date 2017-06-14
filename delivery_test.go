package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	//"strings"
)

type fmturltestpair struct {
	data   postback
	result string
}

//redigo mock interface
type redisMock struct {
	mock.Mock
}

func (o redisMock) Do(cmd string, a ...interface{}) (interface{}, error) {

	args := o.Called(cmd, a)

	if args[0] == nil {
		return nil, nil
	}

	return args.Get(0).(interface{}), args.Error(1)
}

func (o redisMock) Close() error {

	args := o.Called()
	return args.Error(0)
}

func (o redisMock) Err() error {

	args := o.Called()
	return args.Error(0)
}

func (o redisMock) Flush() error {

	args := o.Called()
	return args.Error(0)
}

func (o redisMock) Receive() (interface{}, error) {

	args := o.Called()
	return args.Get(0).([]uint8), args.Error(1)
}

func (o redisMock) Send(a string, i ...interface{}) error {

	args := o.Called(a, i)
	return args.Error(0)
}

var urlTests = []fmturltestpair{
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={money}", Data: map[string]string{"money": "100 dollars"}}, result: "https://httpbin.org/get?evil=100+dollars"}, //url + key + data
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={money}"}, result: "https://httpbin.org/get?evil="},                                                             //url + key + no data
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil="}, result: "https://httpbin.org/get?evil="},                                                                    //url + no key + no data
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil=", Data: map[string]string{"money": "100 dollars"}}, result: "https://httpbin.org/get?evil="},                   //url + no key + data
	{data: postback{Method: "GET", Url: "{money}", Data: map[string]string{"money": "100 dollars"}}, result: "100+dollars"},                                                           //no url + key + data
	{data: postback{Method: "GET", Url: ""}, result: ""},                                                                                                                              //no url + no key + no data
	{data: postback{Method: "GET", Url: "", Data: map[string]string{"money": "100 dollars"}}, result: ""},                                                                             //no url + no key + data
	{data: postback{Method: "GET", Url: "{money}"}, result: ""},                                                                                                                       //no url + key + no data
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={money}&other={other}", Data: map[string]string{"money": "100 dollars", "other": "something else"}}, result: "https://httpbin.org/get?evil=100+dollars&other=something+else"},                 //url + mult key + mult data
	{data: postback{Method: "GET", Url: "{money}&other={other}", Data: map[string]string{"money": "100 dollars", "other": "something else"}}, result: "100+dollars&other=something+else"},                                                                           //no url + mult key + mult data
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={money}&other={other}"}, result: "https://httpbin.org/get?evil=&other="},                                                                                                                      //url + mult key + no data
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={money}&other={other}", Data: map[string]string{"other": "something else"}}, result: "https://httpbin.org/get?evil=&other=something+else"},                                                    //url + mult key + incomplete data at end
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={money}&other={other}", Data: map[string]string{"money": "100 dollars"}}, result: "https://httpbin.org/get?evil=100+dollars&other="},                                                          //url + mult key + incomplete data at beg
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={money}", Data: map[string]string{"other": "something else"}}, result: "https://httpbin.org/get?evil="},                                                                                       //url + key + different data
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={other}", Data: map[string]string{"other": "! @ # $ % ^ & * ( ) + ? > < , ; :"}}, result: "https://httpbin.org/get?evil=%21+%40+%23+%24+%25+%5E+%26+%2A+%28+%29+%2B+%3F+%3E+%3C+%2C+%3B+%3A"}, //check urlencoding
}

var TestFormatUrlRegexKeys = []fmturltestpair{
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={$money}", Data: map[string]string{"$money": "100 dollars"}}, result: "https://httpbin.org/get?evil=100+dollars"}, //url + key with regexp + data with regexp
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={$money}", Data: map[string]string{"money": "100 dollars"}}, result: "https://httpbin.org/get?evil="},             //url + key with regexp + data without regexp
	{data: postback{Method: "GET", Url: "{.*}", Data: map[string]string{".*": "100 dollars", "other": "something else"}}, result: "100+dollars"},                                        //key is a regex
	{data: postback{Method: "GET", Url: "{.*?$[]^$}", Data: map[string]string{".*?$[]^$": "100 dollars", "other": "something else"}}, result: "100+dollars"},                            //key has all regexp symbols
	{data: postback{Method: "GET", Url: "{{}}", Data: map[string]string{"{}": "100 dollars", "other": "something else"}}, result: "100+dollars"},                                        //key is a regex
	{data: postback{Method: "GET", Url: "https://httpbin.org/get?evil={$money}"}, result: "https://httpbin.org/get?evil="},                                                              //key contains regexp no data
	{data: postback{Method: "GET", Url: "{$money*}{other?}", Data: map[string]string{"$money*": "100 dollars", "other?": "something else"}}, result: "100+dollarssomething+else"},
}

func TestFormatUrl(t *testing.T) {
	for _, pair := range urlTests {
		r := formatUrl(pair.data)
		assert.Equal(t, pair.result, r.Url, "formatUrl() didn't return the expected formatted url.")
	}
}

func TestFormatUrlRegex(t *testing.T) {

	for _, pair := range TestFormatUrlRegexKeys {
		r := formatUrl(pair.data)
		assert.Equal(t, pair.result, r.Url, "formatUrl() didn't return the expected formatted url.")
	}
}

func TestSendRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "some string")
	}))
	defer testServer.Close()

	testURl := testServer.URL

	resp, err := sendRequest(testURl, "GET")

	assert.Equal(t, nil, err, "err should be nil")
	assert.Equal(t, "200", resp.responseCode, "sendRequest() didn't return the expected status code.")
	assert.Equal(t, "some string\n", resp.responseBody, "sendRequest() didn't return the expected response body.")

	resp, err = sendRequest(testURl, "POST")
	assert.NotEqual(t, nil, err, "err shoudln't be nil")
	assert.Equal(t, errors.New("error, only GET requests are supported"), err)

}

func TestProcess(t *testing.T) {
	client := new(redisMock)
	client.On("Do", "RPOP", []interface{}{"data"}).Return([]uint8("{\"method\":\"GET\",\"url\":\"https:\\/\\/httpbin.org\\/get?evil={$money}\",\"data\":{\"$money\":\"100 dollars\"}}"), nil)
	obj, _ := process(client)

	url := "https://httpbin.org/get?evil={$money}"
	data := "100 dollars"

	assert.Equal(t, "GET", obj.Method, "Incorrect Http method")
	assert.Equal(t, url, obj.Url, "Url does not match")
	assert.Equal(t, data, obj.Data["$money"], "Given data does not match")
	assert.Equal(t, 1, len(obj.Data), "There should only be one data key")
}

func TestProcessEmptyRequest(t *testing.T) {
	client := new(redisMock)
	client.On("Do", "RPOP", []interface{}{"data"}).Return(nil, nil)
	obj, _ := process(client)

	assert.Nil(t, obj, "obj is supposed to be nil")

}
