httpjc
========

A simple HTTP server that listens on port 8080 to provide password encryption
service. The hash is computed as base64 encoded string of the SHA512 hash
of the provided password string.

## Prerequisites
* [Go tools](https://golang.org/doc/install)
* [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Installation
```bash
go get github.com/edchen1025/httpjc
```

## Launch the Server
```bash
# Start standalone httpjc server
httpjc &
```

## POST with a Password
```bash
# Submit a password encryption request using /hash endpoint
curl -d "password=angryMoney" http://localhost:8080/hash
# An identifier will be returned
1
```

## GET a Hash
```bash
# Get a password hash using /hash/# endpoint
curl http://localhost:8080/hash/1
# A quoted string will be returned
"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
```

## Statistics
```bash
# Request statistical information using /stats endpoint
curl http://localhost:8080/stats
# A JSON object will be returned. The "total" key has the value of the total
# number of requests; the "average" key has the value of the average time taken
# of each request in microseconds.
{"total":1,"average":12}
```

## Shutdown
```bash
# Gracefully shutdown this server using /shutdown endpoint
curl http://localhost:8080/shutdown
```
