package main

import (
  "fmt"
  "github.com/garyburd/redigo/redis"
  "time"
  "encoding/json"
  "regexp"
//  "net/http"
)


type postback struct {
  Method string `json:"method"`
  Url string `json:"url"`
  Data map[string]string
}

func main() {

  fmt.Println("Delivery server started \nPress Ctrl-C to quit.")


  client, err := redis.Dial("tcp", ":6369")
  if err != nil {
    panic(err)
  }
  defer client.Close()

  for {
    request, err := client.Do("RPOP", "data")

    if request != nil {
      request, _ := redis.String(request, err)
      temp := postback{}

      if err := json.Unmarshal([]byte(request), &temp); err != nil {
        panic(err)
      }
      for key, value := range temp.Data {
        fmt.Println("key: " + key + " value: " + value)
      }
      fmt.Println(temp)
    } else {
      fmt.Println("the queue is empty")
    }

    time.Sleep(1000 * time.Millisecond)
  }
}




