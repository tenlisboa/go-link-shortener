# The GO link shortener

This is a project with the focus of studying the implementation of observability principles, using OpenTelemetry framework
the prime focus is implement a very small service and trace this activity using `opentelemetry-go` and `zipkin`.

This is provisioned with `docker-compose` and everyone can run the project, the requirement is `docker` and `docker-compose`.

## The service
The structure is quite simple.

`/short` receives a JSON with one parameter `url` that is the link you want to shorten, and it returns
another JSON with `url` that is the link shortened.

Clicking on the link provided by `/short` you will get into the `/{hash}` endpoint, that will redirect you to the link associated 
previously.

## Tech Stack
- Go
- Mux
- Redis
- OpenTelemetry
- Zipkin
- Docker