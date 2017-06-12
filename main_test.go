package main

import (
  "testing"
  "strings"
)

type testpair struct {
  data postback
  result string
}


func TestFormatUrl(t *testing.T) {
  pair := testpair{data: postback{Method:"GET", Url: "https://httpbin.org/get?evil={$money}", Data: map[string]string{"$money": "100 dollars"}}, result:"https://httpbin.org/get?evil=100+dollars"}

  r := formatUrl(pair.data)
  if (strings.Compare(r.Url, pair.result) != 0) {
    t.Errorf("The url wasn't formatted correctly. was expecting %s got %s", pair.result, r.Url)
  }

  pair = testpair{data: postback{Method:"GET", Url: "https://httpbin.org/get?evil={$money}"}, result:"https://httpbin.org/get?evil="}

  r = formatUrl(pair.data)
  if (strings.Compare(r.Url, pair.result) != 0) {
    t.Errorf("The url wasn't formatted correctly. was expecting %s got %s", pair.result, r.Url)
  }

  pair = testpair{data: postback{Method:"GET", Url: "https://httpbin.org/get?evil="}, result:"https://httpbin.org/get?evil="}

  r = formatUrl(pair.data)
  if (strings.Compare(r.Url, pair.result) != 0) {
    t.Errorf("The url wasn't formatted correctly. was expecting %s got %s", pair.result, r.Url)
  }

  pair = testpair{data: postback{Method:"GET", Url: "", Data: map[string]string{"$money": "100 dollars"}}, result:""}

  r = formatUrl(pair.data)
  if (strings.Compare(r.Url, pair.result) != 0) {
    t.Errorf("The url wasn't formatted correctly. was expecting %s got %s", pair.result, r.Url)
  }
  pair = testpair{data: postback{Method:"GET", Url: "https://httpbin.org/get?evil="}, result:"https://httpbin.org/get?evil="}


  r = formatUrl(pair.data)
  if (strings.Compare(r.Url, pair.result) != 0) {
    t.Errorf("The url wasn't formatted correctly. was expecting %s got %s", pair.result, r.Url)
  }
  pair = testpair{data: postback{Method:"GET", Url: "https://httpbin.org/get?evil={$money}&other={other}", Data: map[string]string{"$money": "100 dollars", "other":"something else"}}, result:"https://httpbin.org/get?evil=100+dollars&other=something+else"}


  r = formatUrl(pair.data)
  if (strings.Compare(r.Url, pair.result) != 0) {
    t.Errorf("The url wasn't formatted correctly. was expecting %s got %s", pair.result, r.Url)
  }
  pair = testpair{data: postback{Method:"GET", Url: "https://httpbin.org/get?evil={$money}&other={other}"}, result:"https://httpbin.org/get?evil=&other="}


  r = formatUrl(pair.data)
  if (strings.Compare(r.Url, pair.result) != 0) {
    t.Errorf("The url wasn't formatted correctly. was expecting %s got %s", pair.result, r.Url)
  }
  pair = testpair{data: postback{Method:"GET", Url: "https://httpbin.org/get?evil={$money}&other={other}", Data: map[string]string{"other":"something else"}}, result:"https://httpbin.org/get?evil=&other=something+else"}


  r = formatUrl(pair.data)
  if (strings.Compare(r.Url, pair.result) != 0) {
    t.Errorf("The url wasn't formatted correctly. was expecting %s got %s", pair.result, r.Url)
  }
  pair = testpair{data: postback{Method:"GET", Url: "{.*}", Data: map[string]string{".*": "100 dollars", "other":"something else"}}, result:"100+dollars"}


  r = formatUrl(pair.data)
  if (strings.Compare(r.Url, pair.result) != 0) {
    t.Errorf("The url wasn't formatted correctly. was expecting %s got %s", pair.result, r.Url)
  }
  pair = testpair{data: postback{Method:"GET", Url: "{$money*}{other?}", Data: map[string]string{"$money*": "100 dollars", "other?":"something else"}}, result:"100+dollarssomething+else"}


  r = formatUrl(pair.data)
  if (strings.Compare(r.Url, pair.result) != 0) {
    t.Errorf("The url wasn't formatted correctly. was expecting %s got %s", pair.result, r.Url)
  }

}
