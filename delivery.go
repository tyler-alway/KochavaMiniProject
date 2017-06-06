package main

import (
  "fmt"
  "strings"
  "github.com/garyburd/redigo/redis"
  "time"
  "encoding/json"
  "regexp"
  "net/url"
  "net/http"
  "log"
  "os"
  "errors"
)

//struct to hold the postback object after beign parsed from json
type postback struct {
  Method string `json:"method"`
  Url string `json:"url"`
  Data map[string]string
}

func main() {
  fmt.Println("Delivery server started \nPress Ctrl-C to quit.")


  //Opens file for error logging
  f, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  log.SetOutput(f)


  //opens a tcp connection on port 6369 to the redis-server (queue)
  client, err := redis.Dial("tcp", ":6369")
  if err != nil {
    //This will be logged to the log.txt file and should exit the program 
    log.Println(err)
    panic(err)
  }

  //Authenticates to connect to the resis server
  _, err = client.Do("AUTH", "anotherPassword")

  //if the connection to the server fails panic
  if err != nil {
    //This will be logged to the log.txt file and should exit the program 
    log.Println(err)
    panic(err)
  }

  //Dont allow the connection to the redis server to be closed utill the main function as finished executing
  defer client.Close()

  //for loop to continue to pull postback objects out of the queue and push them to the endpoint
  for {
    //pulls a postback object off the queue
    request, err := client.Do("RPOP", "data")
    if (err != nil) {
      log.Println(err)
    }

    //if there is a postback object in the queue
    //(if the queue was empty request will be nil)
    if request != nil {
      request, _ := redis.String(request, err)
      //creates a new postback object (struct) for the json string to be parsed into 
      temp := postback{}

      //parses the json string into the postback object 
      //if the object is not valid json panic
      if err := json.Unmarshal([]byte(request), &temp); err != nil {
        log.Println(err)
        continue
      }

      temp = formatUrl(temp)

      err = sendRequest(temp.Url, temp.Method)
      //If there is an error log it
      if(err != nil) {
        log.Println(err)
      }

    //else the queue is empty so try and remove another object
    } else {
      fmt.Println("the queue is empty")
    }
    //TODO this needs to be removed it is here so that the testing values are readable (gives me time to read the testing outputs)
    time.Sleep(1000 * time.Millisecond)
  }
  //close the log file
  defer f.Close()

}


/*
* Name: formatUrl
* Description: Function format the given postback object 
* Parameters: Takes in a postback object 
* Returns: The formatted data obj
*/
func formatUrl(data postback) postback {
  //loop though the data section of the postback object replace  {xxx} with Date[xxx]
  for key, value := range data.Data {
    value = url.QueryEscape(value)
    key = "{" + key +"}"
    re := regexp.MustCompile(key)
    data.Url = re.ReplaceAllString(data.Url, value)
  }

  //if there are any unmatched {xxx} strings remove them from the final url
  re := regexp.MustCompile("{.*}")
  data.Url = re.ReplaceAllString(data.Url, "")

  return data
}


/*
* Name: sendRequest
* Description: Function to send get requests to the endpoint
* Parameters: requestType and the pre formatted url to be sent
* Returns: None
*/
func sendRequest(url string, requestType string) error {

  //Make sure that the request type is GET 
  //Only GET request are supported
  if (strings.Compare("GET", requestType) != 0) {
    return errors.New("error, only GET requests are supported")
  }
  //Send the GET request
  t1 := time.Now()
  response, err := http.Get(url)
  t2 := time.Now()
  //Finds the time it took to send the request and recive the response
  d := t2.Sub(t1)
  log.Println(d)

  if err != nil {
    return err
  } else {
    defer response.Body.Close()
    log.Println(response.StatusCode)
    log.Println(response.Body)
  }
  return nil
}



