FROM golang:1.19.2-alpine3.16 as builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o asker ./cmd/asker

FROM alpine:3.16 as runner

WORKDIR /app

RUN apk update && apk upgrade && apk add bash

COPY --from=builder /src/asker /src/pkg/utils/wait-for-it.sh ./
RUN chmod 777 ./wait-for-it.sh
