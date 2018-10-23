# goms

<!-- Badger start badges -->
[![Status of the build](https://badger.spt-engprod-pro.schibsted.io/badge/travis/Yapo/goms)](https://travis.schibsted.io/Yapo/goms)
[![Testing Coverage](https://badger.spt-engprod-pro.schibsted.io/badge/coverage/Yapo/goms)](https://reports.spt-engprod-pro.schibsted.io/#/Yapo/goms?branch=master&type=push&daterange&daterange)
[![Style/Linting issues](https://badger.spt-engprod-pro.schibsted.io/badge/issues/Yapo/goms)](https://reports.spt-engprod-pro.schibsted.io/#/Yapo/goms?branch=master&type=push&daterange&daterange)
[![Badger](https://badger.spt-engprod-pro.schibsted.io/badge/flaky_tests/Yapo/goms)](https://databulous.spt-engprod-pro.schibsted.io/test/flaky/Yapo/goms)
[![Badger](https://badger.spt-engprod-pro.schibsted.io/badge/quality_index/Yapo/goms)](https://databulous.spt-engprod-pro.schibsted.io/quality/repo/Yapo/goms)
[![Badger](https://badger.spt-engprod-pro.schibsted.io/badge/engprod/Yapo/goms)](https://github.schibsted.io/spt-engprod/badger)
<!-- Badger end badges -->

Goms is the official golang microservice template for Yapo.

## A few rules

* Goms was built following [Clean Architecture](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164) so, please, familiarize yourself with it and let's code great code!

* Goms has great [test coverage](https://quality-gate.schibsted.io/#/Yapo/goms) and [examples](https://github.schibsted.io/Yapo/goms/search?l=Go&q=func+Test&type=&utf8=%E2%9C%93) of how good testing can be done. Please honor the effort and keep your test quality in the top tier.

* Goms is not a silver bullet. If your service clearly doesn't fit in this template, let's have a [conversation](mailto:dev@schibsted.cl)

* [README.md](README.md) is the entrypoint for new users of your service. Keep it up to date and get others to proof-read it.

## How to run the service

* Create the dir: `~/go/src/github.schibsted.io/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.schibsted.io/Yapo
  $ git clone git@github.schibsted.io:Yapo/goms.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd goms
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
    docker-publish       Push docker image to containers.schibsted.io
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
  

## Creating a new service

* Create a repo for your new service on: https://github.schibsted.io/Yapo
* Rename your goms dir to your service name:
  ```
  $ mv goms YourService
  ```
* Update origin: 
  ```
  # https://help.github.com/articles/changing-a-remote-s-url/
  $ git remote set-url origin git@github.schibsted.io:Yapo/YourService.git
  ```

* Replace every goms reference to your service's name:
  ```
  $ git grep -l goms | xargs sed -i.bak 's/goms/yourservice/g'
  $ find . -name "*.bak" | xargs rm
  ```

* Go through the code examples and implement your service
  ```
  $ git grep -il fibonacci
  README.md
  cmd/goms/main.go
  pkg/domain/fibonacci.go
  pkg/domain/fibonacci_test.go
  pkg/interfaces/handlers/fibonacci.go
  pkg/interfaces/handlers/fibonacci_test.go
  pkg/interfaces/loggers/fibonacciInteractorLogger.go
  pkg/interfaces/repository/fibonacci.go
  pkg/interfaces/repository/fibonacci_test.go
  pkg/usecases/getNthFibonacci.go
  pkg/usecases/getNthFibonacci_test.go
  ```

* Enable TravisCI
  - Go to your service's github settings -> Hooks & Services -> Add Service -> Travis CI
  - Fill in the form with the credentials you obtain from https://travis.schibsted.io/profile/
  - Sync your repos and organizations on Travis
  - Make a push on your service
  - The push should trigger a build. If it didn't ensure that it is enabled on the travis service list
  - Enjoy! This should automatically enable quality-gate reports and a few other goodies

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

### GET  /api/v1/fibonacci
Implements the Fibonacci Numbers with Clean Architecture

#### Request
{
	"n": int - Ask for the nth fibonacci number
}

#### Response

```javascript
200 OK
{
	"Result": int - The nth fibonacci number
}
```

#### Error response
```javascript
400 Bad Request
{
	"ErrorMessage": string - Explaining what went wrong
}
```

### Contact
dev@schibsted.cl
