# Postback Miniproject

## Description
##### Kochava Miniproject :: "Postback Delivery"
This project takes http requests into ingest.php and stores the POST data sent in into a Redis server (postback Delivery Queue) then a go Delivery agent pulls requests out of the Delivery Queue, formats them, and sends them to the endpont location. Error and Delivery response data will be logged to log.txt.
## Setup
Required software (must be installed):
- NGINX
- PHP 7
- Go
- Redis
- Predis (PHP Redis client)
- Redigo (Go Redis client)

### Notes:
##### NGINX:
Useful resouce for installing and getting NGINX working with PHP:
https://www.digitalocean.com/community/tutorials/how-to-install-linux-nginx-mysql-php-lemp-stack-in-ubuntu-16-04
Move the PHP script so NGINX will run it  `cp ingest.php /var/www/html/ingest.php`
##### Go
Go install instructions
https://golang.org/doc/install
You will also need to set you $GOPATH env variable:
https://astaxie.gitbooks.io/build-web-application-with-golang/en/01.2.html

##### Redis
Installing Redis
https://redis.io/topics/quickstart
You will also need to edit the redis.conf file located at `/etc/redis/redis.conf`
You will change the lines labeled
`port <Your Port>`
and
`requirepass <Your Redis Password>`

##### Predis
Nrk's PEAR channel:
http://pear.nrk.io/

##### Redigo
https://github.com/garyburd/redigo
You will also need $GOPATH to be set

##### ingest.php
Change line 18 to: `'password' => <Your Redis Password>`
##### delivery.go
Cnage line 42 to: `_, err = client.Do("AUTH", "<Tour Redis Password>")`

## Usage
To run the Redis server: `redis-server /etc/redis/redis.conf`
To have NGINX run the PHP ingestor (ingest.php): `cp ingest.php /var/www/html/ingest.php`
To run the Go Delivery Agnet: `go run delivery.go`
Or to run go in the background
1) `go build delivery.go`
2) `./delivery &`
###### Sample curls:
`$ curl -X POST -H "Content-Type: application/json" -d '{"endpoint":{"method":"GET","url":"http://localhost:80/ingest.php?title={mascot}&image={location}&foo={bar}"},"data":[{"mascot":"Gopher","location":"https://blog.golang.org/gopher/gopher.png"}]}'  159.203.164.144/ingest.php`

`$ curl -X POST -H "Content-Type: application/json" -d '{"endpoint":{"method":"GET","url":"http://sample_domain_endpoint.com/data?title={mascot}&image={location}&foo={bar}"},"data":[{"mascot":"Gopher","location":"https://blog.golang.org/gopher/gopher.png"}]}'  159.203.164.144/ingest.php`
