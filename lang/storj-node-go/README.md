# storj-node
> implemented in Go

This is an experiment to implement a storj node in Go and explore the language features, possible frameworks, persistence layers, and the general Go ecosystem for use in our next iteration.

## Goals

- Spec out the full skeleton of API endpoints as per the [Spec draft](../../docs/spec-draft.md)
- Package a binary
- Build and dockerize the server and create a workflow for it
- Document any pros and cons and our experience with Go overall

## Development with single docker image

You can run a Docker container on local

```bash
# Build a development container 

docker build -t storj-node-dev -f dockerfiles/dev.dockerfile

# Run the container 

docker run -v ~/path/to/storj-node-go/code/:/go/src/storj-node-go -p 8080:8080 -d <container-id>

# Docker exec into the container 

docker exec -it gifted_heisenberg /bin/sh

```
When inside the container, run `go run main.go` and it'll run the Go server

## Development with Docker-compose
*WIP*
