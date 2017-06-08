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
##### NGINX:
Useful resource for installing and getting NGINX working with PHP:
https://www.digitalocean.com/community/tutorials/how-to-install-linux-nginx-mysql-php-lemp-stack-in-ubuntu-16-04

Add local environment variables to php-fpm by adding these lines:
`env["REDISPORT"] = 6369`
`env["REDISPASS"] = anotherPassword`
to the end of `vim /etc/php/7.0/fpm/pool.d/www.conf`

Then restart php-fpm: `service php7.0-fpm restart`

Move the PHP script so NGINX will run it:  `cp ingest.php /var/www/html/ingest.php`
##### Go
Go install instructions:
https://golang.org/doc/install
You will also need to set you $GOPATH env variable:
https://astaxie.gitbooks.io/build-web-application-with-golang/en/01.2.html

##### Redis
Installing Redis:
https://redis.io/topics/quickstart
You will also need to edit the redis.conf file located at `/etc/redis/redis.conf`
You will change the lines labeled
`port <Your Redis Port>`
and
`requirepass <Your Redis Password>`

Then run these commands:
`export REDISPORT=<Your Redis Port>`
`export REDISPASS=<Your Redis Password>`

##### Predis
Nrk's PEAR channel:
http://pear.nrk.io/

##### Redigo
https://github.com/garyburd/redigo
You will also need $GOPATH to be set


## Usage
To run the Redis server: `redis-server /etc/redis/redis.conf`
To have NGINX run the PHP ingest agent (ingest.php): `cp ingest.php /var/www/html/ingest.php`
To run the Go Delivery Agent: `go run delivery.go`
Or to run go in the background
1) `go build delivery.go`
2) `./delivery &`

###### Sample curls:
`$ curl -X POST -H "Content-Type: application/json" -d '{"endpoint":{"method":"GET","url":"http://localhost/ingest.php?title={mascot}&image={location}&foo={bar}"},"data":[{"mascot":"Gopher","location":"https://blog.golang.org/gopher/gopher.png"}]}'  http://localhost/ingest.php`

`$ curl -X POST -H "Content-Type: application/json" -d '{"endpoint":{"method":"GET","url":"http://sample_domain_endpoint.com/data?title={mascot}&image={location}&foo={bar}"},"data":[{"mascot":"Gopher","location":"https://blog.golang.org/gopher/gopher.png"}]}'  http://localhost/ingest.php`
