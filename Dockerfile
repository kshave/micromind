#syntax=docker/dockerfile:1

## Build
FROM golang:buster AS Build

WORKDIR /service

COPY . ./
RUN go mod download

RUN go build -o /micromind

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=BUILD /micromind /micromind

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/micromind"]
