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
	"strings"
	"time"
)

//struct to hold the postback object after beign parsed from json
type postback struct {
	Method string `json:"method"`
	Url    string `json:"url"`
	Data   map[string]string
}

func main() {

	//Opens file for error logging
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)

	redisPort := os.Getenv("REDISPORT")
	redisPass := os.Getenv("REDISPASS")

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
		//pulls a postback object off the queue
		request, err := client.Do("RPOP", "data")
		if err != nil {
			log.Println(err)
		} else if request != nil {
			request, _ := redis.String(request, err)
			//creates a new postback object (struct) for the json string to be parsed into
			temp := postback{}

			//parses the json string into the postback object
			if err := json.Unmarshal([]byte(request), &temp); err != nil {
				log.Println(err)
				continue
			}

			temp = formatUrl(temp)

			err = sendRequest(temp.Url, temp.Method)
			if err != nil {
				log.Println(err)
			}
		}
	}

	defer client.Close()
	defer f.Close()

}

//Name: formatUrl
//Description: Function format the given postback object
//Parameters: Takes in a postback object
//Returns: The formatted data obj
func formatUrl(data postback) postback {
	//loop though the data section of the postback object replace  {xxx} with Date[xxx]
	for key, value := range data.Data {
		value = url.QueryEscape(value)
		key = "{" + key + "}"
		re := regexp.MustCompile(key)
		data.Url = re.ReplaceAllString(data.Url, value)
	}

	//if there are any unmatched {xxx} strings remove them from the final url
	re := regexp.MustCompile("{.*}")
	data.Url = re.ReplaceAllString(data.Url, "")

	return data
}

//Name: sendRequest
//Description: Function to send get requests to the endpoint
//Parameters: requestType and the pre formatted url to be sent
//Returns: None
func sendRequest(url string, requestType string) error {

	if strings.Compare("GET", requestType) != 0 {
		return errors.New("error, only GET requests are supported")
	}

	t1 := time.Now()
	response, err := http.Get(url)
	t2 := time.Now()

	d := t2.Sub(t1)
	log.Println(d)

	if err != nil {
		return err
	} else {
		defer response.Body.Close()
		log.Println(response.StatusCode)
		bs, _ := ioutil.ReadAll(response.Body)
		log.Println(string(bs))
	}
	return nil
}
