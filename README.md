# Postback Miniproject

## Description
##### Kochava Miniproject :: "Postback Delivery"
This project accepts http requests in ingest.php, stores the POST data into a Redis server (Delivery Queue), and a golang Delivery agent pulls requests out of the Delivery Queue, formats them, and sends them to the endpoint location. Error and Delivery response data will be logged to log.txt.
## Setup
Required software (must be installed):
- NGINX
- PHP 7
- Go
- Redis
- Predis (PHP Redis client)
- Redigo (Go Redis client)

### Notes:

##### Redis
Installing Redis:
https://redis.io/topics/quickstart
You will also need to edit the redis.conf file located at `/etc/redis/redis.conf`
You will change the lines labeled
`port <Your Redis Port>`
and
`requirepass <Your Redis Password>`

Then run these commands:<br/>
`export REDISPORT=<Your Redis Port>` <br/>
`export REDISPASS=<Your Redis Password>`


##### NGINX:
Useful resource for installing and getting NGINX working with PHP:
https://www.digitalocean.com/community/tutorials/how-to-install-linux-nginx-mysql-php-lemp-stack-in-ubuntu-16-04

Open the config.php file and replace `<Your Redis Port>` with your redis port
and replace `<Your Redis Passowrd>` with your redis password

Then move the PHP script and config to NGINX:  
`cp ingest.php /var/www/html/ingest.php`<br/>
`cp config.php /var/www/html/config.php`


##### Go
Go install instructions:
https://golang.org/doc/install

You will also need to set you $GOPATH env variable:
https://astaxie.gitbooks.io/build-web-application-with-golang/en/01.2.html

##### Predis
Nrk's PEAR channel:
http://pear.nrk.io/

##### Redigo
https://github.com/garyburd/redigo (Note: You will need $GOPATH to be set)


## Tests
For running test you will need get the Testify package from https://github.com/stretchr/testify
Or to install the newest version run `go get github.com/stretchr/testify`

To run tests: `go test`


## Usage
To run the Redis server: `redis-server /etc/redis/redis.conf`
To have NGINX run the PHP ingest agent (ingest.php): `cp ingest.php /var/www/html/ingest.php`
To run the Go Delivery Agent: `go run delivery.go`
Or to run go in the background
1) `go build delivery.go`
2) `./delivery &`

##### Sample curls:

`$ curl -X POST -H "Content-Type: application/json" -d '{"endpoint":{"method":"GET","url":"http://sample_domain_endpoint.com/data?title={mascot}&image={location}&foo={bar}"},"data":[{"mascot":"Gopher","location":"https://blog.golang.org/gopher/gopher.png"}]}'  http://localhost/ingest.php`

###### Expected HTTP Request:

`GET http://sample_domain_endpoint.com/data?title=Gopher&image=https%3A%2F%2Fblog.golang.org%2Fgopher%2Fgopher.png&foo=`

`$ curl -X POST -H "Content-Type: application/json" -d '{"endpoint":{"method":"GET","url":"http://localhost/ingest.php?title={mascot}&image={location}&foo={bar}"},"data":[{}]}'  http://localhost/ingest.php`

###### Expected HTTP Request:

`GET http://localhost/ingest.php?title=&image=&foo=`

`$ curl -X POST -H "Content-Type: application/json" -d '{"endpoint":{"method":"GET","url":"https://httpbin.org/get?evil={$money}"},"data":[{"$money": "100 dollars"}]}'  localhost/ingest.php`

###### Expected HTTP Request:

`GET https://httpbin.org/get?evil=100+dollars`
