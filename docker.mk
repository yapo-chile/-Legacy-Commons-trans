## Create docker image based on docker/dockerfile
docker-build:
	@scripts/commands/docker-build.sh

## Push docker image to containers.mpi-internal.com
docker-publish:
	@scripts/commands/docker-publish.sh

## Attach to this service's currently running docker container output stream
docker-attach:
	@scripts/commands/docker-attach.sh

## Start all required docker containers for this service
docker-compose-up:
	@scripts/commands/docker-compose-up.sh

## Stop all running docker containers for this service
docker-compose-down:
	@scripts/commands/docker-compose-down.sh
