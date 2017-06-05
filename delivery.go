package main

import (
  "fmt"
  "github.com/garyburd/redigo/redis"
  "time"
)

func main() {

  fmt.Println("Delivery server started \nPress Ctrl-C to quit.")


  client, err := redis.Dial("tcp", ":6369")
  if err != nil {
    panic(err)
  }
  defer client.Close()

  for {
    pushback, err := client.Do("RPOP", "test")

    if pushback != nil {
      fmt.Println(redis.String(pushback, err))
    } else {
      fmt.Println("the queue is empty")
    }
    time.Sleep(1000 * time.Millisecond)
  }
}




