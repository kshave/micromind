#syntax=docker/dockerfile:1

## Build
FROM golang:buster AS Build

WORKDIR /service

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /micromind

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=BUILD /micromind /micromind

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/micromind"]