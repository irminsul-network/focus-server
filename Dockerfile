FROM golang:1.23-alpine AS build
LABEL authors="boby"

WORKDIR /garage

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY src/ ./src

RUN go build -v -o app ./src

FROM alpine

WORKDIR /car

VOLUME /car/data/

COPY --from=build /garage/app ./

EXPOSE 8080


ENTRYPOINT ["/car/app"]



