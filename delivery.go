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
	defer f.Close()
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
	defer client.Close()

	//Authenticates to connect to the resis server
	_, err = client.Do("AUTH", redisPass)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	//for loop to continue to pull postback objects out of the queue and push them to the endpoint
	for {

		redisObj, err := fetchPostbackObj(client)
		if err != nil {
			log.Println(err)
		} else if redisObj != nil {
			obj := formatUrl(*redisObj)
			response, err := sendRequest(obj.Url, obj.Method)
			if err != nil {
				log.Println(err)
			} else if response != nil {
				log.Println(response.responseCode)
				log.Println(response.responseTime)
				log.Println(response.responseBody)
			}
		}
	}
}

//Name: fetchPostbackObj(client Conn)
//Description: This function fetches a postback obj from redis
//Parameters: Takes in a redis connection
//Returns: responseData obj, error
func fetchPostbackObj(client redis.Conn) (*postback, error) {
	//pulls a postback object off the queue
	str, err := redis.String(client.Do("RPOP", "data"))
	if err != nil {
		return nil, nil
	}

	obj := postback{}
	//parses the json string into the postback object
	if err := json.Unmarshal([]byte(str), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

//Name: formatUrl
//Description: Function format the given postback object
//Parameters: Takes in a postback object
//Returns: The formatted data obj
func formatUrl(data postback) postback {
	//loop though the data section of the postback object replace {xxx} with Date[xxx]
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
