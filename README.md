# Postback Miniproject

## Description
##### Kochava Miniproject :: "Postback Delivery"
This project accepts http requests in ingest.php, stores the POST data into a Redis server (Delivery Queue), and a golang Delivery agent pulls requests out of the Delivery Queue, formats them, and sends them to the endpoint location. Error and Delivery response data will be logged to log.txt.

## Installation

Docker must be installed:

https://docs.docker.com/engine/installation/

 Clone project directory pointed to be your $GOPATH

 `git clone https://github.com/tyler-alway/KochavaMiniProject.git <your directory>`


## Usage


To run use `docker-compose up` <br/>
To run in the background use `docker-compose up -d` <br/>


## Tests
To run tests use `docker-compose -f docker-compose.tests.yml up`

## Sample curls:

###### Sample curl:

`$ curl -X POST -H "Content-Type: application/json" -d '{"endpoint":{"method":"GET","url":"http://sample_domain_endpoint.com/data?title={mascot}&image={location}&foo={bar}"},"data":[{"mascot":"Gopher","location":"https://blog.golang.org/gopher/gopher.png"}]}'  http://localhost/ingest.php`

###### Expected HTTP Request:

`GET http://sample_domain_endpoint.com/data?title=Gopher&image=https%3A%2F%2Fblog.golang.org%2Fgopher%2Fgopher.png&foo=`

---

###### Sample curl:

`$ curl -X POST -H "Content-Type: application/json" -d '{"endpoint":{"method":"GET","url":"http://localhost/ingest.php?title={mascot}&image={location}&foo={bar}"},"data":[{}]}'  http://localhost/ingest.php`

###### Expected HTTP Request:

`GET http://localhost/ingest.php?title=&image=&foo=`

---

###### Sample curl:

`$ curl -X POST -H "Content-Type: application/json" -d '{"endpoint":{"method":"GET","url":"https://httpbin.org/get?evil={$money}"},"data":[{"$money": "100 dollars"}]}'  localhost/ingest.php`

###### Expected HTTP Request:

`GET https://httpbin.org/get?evil=100+dollars`
