package main

import (
  "fmt"
  "github.com/garyburd/redigo/redis"
  "time"
  "encoding/json"
  "regexp"
//  "net/http"
)

//struct to hold the postback object after beign parsed from json
type postback struct {
  Method string `json:"method"`
  Url string `json:"url"`
  Data map[string]string
}

func main() {
  fmt.Println("Delivery server started \nPress Ctrl-C to quit.")

  //opens a tcp connection on port 6369 to the redis-server (queue)
  client, err := redis.Dial("tcp", ":6369")
  //if the connection to the server fails panic
  if err != nil {
    panic(err)
  }
  //Dont allow the connection to the redis server to be closed utill the main function as finished executing
  defer client.Close()

  //for loop to continue to pull postback objects out of the queue and push them to the endpoint
  for {
    //pulls a postback object off the queue
    request, err := client.Do("RPOP", "data")

    //if there is a postback object in the queue
    //(if the queue was empty request will be nil)
    if request != nil {
      request, _ := redis.String(request, err)i
      //creates a new postback object (struct) for the json string to be parsed into 
      temp := postback{}

      //parses the json string into the postback object 
      //if the object is not valid json panic
      if err := json.Unmarshal([]byte(request), &temp); err != nil {
        panic(err)
      }

      //loop though the data section of the postback object replace  {xxx} with Date[xxx]
      for key, value := range temp.Data {
        key = "{" + key +"}"
        re := regexp.MustCompile(key)
        temp.Url = re.ReplaceAllString(temp.Url, value)
      }

      //if there are any unmatched {xxx} strings remove them from the final url
      re := regexp.MustCompile("{.*}")
      temp.Url = re.ReplaceAllString(temp.Url, "")
      fmt.Println(temp.Method + " " + temp.Url)
    //else the queue is empty so try and remove another object
    } else {
      fmt.Println("the queue is empty")
    }
    //TODO this needs to be removed it is here so that the testing values are readable (gives me time to read the testing outputs)
    time.Sleep(1000 * time.Millisecond)
  }
}

