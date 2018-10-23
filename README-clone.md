# __SERVICE__

<!-- Badger start badges -->
[![Status of the build](https://badger.spt-engprod-pro.schibsted.io/badge/travis/Yapo/__SERVICE__)](https://travis.schibsted.io/Yapo/__SERVICE__)
[![Testing Coverage](https://badger.spt-engprod-pro.schibsted.io/badge/coverage/Yapo/__SERVICE__)](https://reports.spt-engprod-pro.schibsted.io/#/Yapo/__SERVICE__?branch=master&type=push&daterange&daterange)
[![Style/Linting issues](https://badger.spt-engprod-pro.schibsted.io/badge/issues/Yapo/__SERVICE__)](https://reports.spt-engprod-pro.schibsted.io/#/Yapo/__SERVICE__?branch=master&type=push&daterange&daterange)
[![Badger](https://badger.spt-engprod-pro.schibsted.io/badge/flaky_tests/Yapo/__SERVICE__)](https://databulous.spt-engprod-pro.schibsted.io/test/flaky/Yapo/__SERVICE__)
[![Badger](https://badger.spt-engprod-pro.schibsted.io/badge/quality_index/Yapo/__SERVICE__)](https://databulous.spt-engprod-pro.schibsted.io/quality/repo/Yapo/__SERVICE__)
[![Badger](https://badger.spt-engprod-pro.schibsted.io/badge/engprod/Yapo/__SERVICE__)](https://github.schibsted.io/spt-engprod/badger)
<!-- Badger end badges -->

__SERVICE__ needs a description here.

## Checklist: Is my service ready?

* [ ] Configure your github repository
  - Open https://github.schibsted.io/Yapo/__SERVICE__/settings
  - Features: Wikis, Restrict editing, Issues, Projects
  - Merge button: Only allow merge commits
  - GitHub Pages: master branch / docs folder
  - Open https://github.schibsted.io/Yapo/goms/settings/branches
  - Default branch: master
  - Protected branches: choose master
  - Protect this branch
    + Require pull request reviews
      - Dismiss stale pull request
    + Require status checks before merging
      - Require branches to be up to date
      - Quality gate code analysis
      - Quality gate coverage
      - Travis-ci
    + Include administrators
* [ ] Enable TravisCI
  - Go to your service's github settings -> Hooks & Services -> Add Service -> Travis CI
  - Fill in the form with the credentials you obtain from https://travis.schibsted.io/profile/
  - Sync your repos and organizations on Travis
  - Create a pull request and make a push on it
  - The push should trigger a build. If it didn't, ensure that it is enabled on the travis service list
  - Enjoy! This should automatically enable quality-gate reports and a few other goodies
* [ ] Get your first PR merged
  - Master should be a protected branch, so the only way to get commits there is via pull request
  - Once the travis build is ok, and you got approval merge it back to master
  - This will allow for the broken badges on top of this readme to display correctly
  - Should them not display after some time, please report it
* [ ] Delete this section
  - It's time for me to leave, I've done my part
  - It's time for you to start coding your new service and documenting your endpoints below
  - Seriously, document your endpoints and delete this section

## How to run __SERVICE__

* Create the dir: `~/go/src/github.schibsted.io/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.schibsted.io/Yapo
  $ git clone git@github.schibsted.io:Yapo/__SERVICE__.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd __SERVICE__
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

## Contact
dev@schibsted.cl
