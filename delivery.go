package main

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//struct to hold the postback object after beign parsed from json
type postback struct {
	Method string            `json:"method"`
	Url    string            `json:"url"`
	Data   map[string]string `json:"data"`
}

//struct to hold the response data object after sending a http request
type responseData struct {
	responseCode string
	responseTime string
	responseBody string
}

func main() {

	//Opens file for error logging
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)

	redisPort := os.Getenv("REDISPORT")
	redisPass := os.Getenv("REDISPASS")

	redisPort = ":" + redisPort

	//opens a tcp connection on port 6369 to the redis-server (queue)
	client, err := redis.Dial("tcp", redisPort)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	//Authenticates to connect to the resis server
	_, err = client.Do("AUTH", redisPass)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	//for loop to continue to pull postback objects out of the queue and push them to the endpoint
	for {
		response, err := process(client)
		if err != nil {
			log.Println(err)
		} else if response != nil {
			log.Println(response.responseCode)
			log.Println(response.responseTime)
			log.Println(response.responseBody)
		}
	}

	defer client.Close()
	defer f.Close()

}


//Name: process(client Conn)
//Description: This function fetches a postback obj from redis and makes the request
//Parameters: Takes in a redis connection
//Returns: responseData obj, error
func process(client redis.Conn) (*responseData, error) {
	//pulls a postback object off the queue
	request, err := client.Do("RPOP", "data")
	if err != nil {
		return nil, err
	} else if request == nil {
		return nil, nil
	}

	data, _ := redis.String(request, err)
	if (err != nil) {
		return nil, err
	}
	//creates a new postback object (struct) for the json string to be parsed into
	temp := postback{}
	//parses the json string into the postback object
	if err := json.Unmarshal([]byte(data), &temp); err != nil {
		return nil, err
	}

	temp = formatUrl(temp)

	response, err := sendRequest(temp.Url, temp.Method)
	if err != nil {
		return nil, err
	}
	return response, err
}

//Name: formatUrl
//Description: Function format the given postback object
//Parameters: Takes in a postback object
//Returns: The formatted data obj
func formatUrl(data postback) postback {
	//loop though the data section of the postback object replace  {xxx} with Date[xxx]
	for key, value := range data.Data {
		value = url.QueryEscape(value)

		re := regexp.MustCompile(regexp.QuoteMeta("{" + key + "}"))
		data.Url = re.ReplaceAllString(data.Url, value)
	}

	//if there are any unmatched {xxx} strings remove them from the final url
	re := regexp.MustCompile("{.*?}")
	data.Url = re.ReplaceAllString(data.Url, "")

	return data
}

//Name: sendRequest
//Description: Function to send get requests to the endpoint
//Parameters: requestType and the pre formatted url to be sent
//Returns: responseData obj, error
func sendRequest(url string, requestType string) (*responseData, error) {

	var httpResponseData responseData

	if strings.Compare("GET", requestType) != 0 {
		return nil, errors.New("error, only GET requests are supported")
	}

	t1 := time.Now()
	response, err := http.Get(url)
	t2 := time.Now()

	duration := t2.Sub(t1)
	httpResponseData.responseTime = duration.String()

	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		httpResponseData.responseCode = strconv.Itoa(response.StatusCode)

		body, _ := ioutil.ReadAll(response.Body)
		httpResponseData.responseBody = string(body[:])

		return &httpResponseData, nil
	}
}
