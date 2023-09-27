FROM golang:1.21.1-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache make
RUN make build

EXPOSE 8080

CMD ["./bin/app"]

