# trans

<!-- Badger start badges -->
[![Status of the build](https://badger.spt-engprod-pro.mpi-internal.com/badge/travis/Yapo/trans)](https://travis.mpi-internal.com/Yapo/trans)
[![Testing Coverage](https://badger.spt-engprod-pro.mpi-internal.com/badge/coverage/Yapo/trans)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/trans?branch=master&type=push&daterange&daterange)
[![Style/Linting issues](https://badger.spt-engprod-pro.mpi-internal.com/badge/issues/Yapo/trans)](https://reports.spt-engprod-pro.mpi-internal.com/#/Yapo/trans?branch=master&type=push&daterange&daterange)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/flaky_tests/Yapo/trans)](https://databulous.spt-engprod-pro.mpi-internal.com/test/flaky/Yapo/trans)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/quality_index/Yapo/trans)](https://databulous.spt-engprod-pro.mpi-internal.com/quality/repo/Yapo/trans)
[![Badger](https://badger.spt-engprod-pro.mpi-internal.com/badge/engprod/Yapo/trans)](https://github.mpi-internal.com/spt-engprod/badger)
<!-- Badger end badges -->

This microservice acts as a proxy between other microservices and a Trans server. The params are passed in a JSON body, and it can be configured to limit what commands can be executed.


## How to run trans

* Create the dir: `~/go/src/github.mpi-internal.com/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.mpi-internal.com/Yapo
  $ git clone git@github.mpi-internal.com:Yapo/trans.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd trans
  $ make start
  ```

* To get a list of available commands:

  ```
  $ make help
  Targets:
    test                 Run tests and generate quality reports
    cover                Run tests and output coverage reports
    coverhtml            Run tests and open report on default web browser
    checkstyle           Run gometalinter and output report as text
    setup                Install golang system level dependencies
    build                Compile the code
    run                  Execute the service
    start                Compile and start the service
    fix-format           Run gofmt to reindent source
    info                 Display basic service info
    docker-build         Create docker image based on docker/dockerfile
    docker-publish       Push docker image to containers.mpi-internal.com
    docker-attach        Attach to this service's currently running docker container output stream
    docker-compose-up    Start all required docker containers for this service
    docker-compose-down  Stop all running docker containers for this service
    help                 This help message
  ```

* If you change the code:

  ```
  $ make start
  ```

* How to run the tests

  ```
  $ make [cover|coverhtml]
  ```

* How to check format

  ```
  $ make checkstyle
  ```

## Endpoints
### GET  /api/v1/healthcheck
Reports whether the service is up and ready to respond.

> When implementing a new service, you MUST keep this endpoint
and update it so it replies according to your service status!

#### Request
No request parameters

#### Response
* Status: Ok message, representing service health

```javascript
200 OK
{
	"Status": "OK"
}
```

### POST  /api/v1/execute/{command}
Sends the specified command to a trans server with the given params in the JSON body

#### Request
params: A JSON object, where the fields are the name of trans params (lowercase), and the values are the values required
by the trans command
```javascript
{
	"params":{
		"param1":"value1",
		...
	}
}
```

#### Response

```javascript
200 OK
{
	"status": "TRANS_OK"
	"response" - A JSON field containing all the values returned by the trans command
	
}
```

#### Error responses
```javascript
400 Bad Request
{
	"status": "TRANS_ERROR"
	"response": {
		"error" - An error message
	}
}
```

```javascript
500 Internal Server Error
{
	"ErrorMessage" - An error message
}
```

### Contact
dev@schibsted.cl
