FROM golang:1.22.3 AS development
WORKDIR /cart
ARG CARTAPP_TOPORT
COPY ./cart .
RUN go mod download
RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      protobuf-compiler
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install github.com/gojuno/minimock/v3/cmd/minimock@latest
RUN go install github.com/cespare/reflex@latest
EXPOSE ${CARTAPP_TOPORT?cartapp_toport_notset}
CMD reflex -r '(\.go$|go\.mod)' make run-all --start-service

FROM golang:1.22.3 AS e2e
WORKDIR /cart
COPY ./cart .


FROM golang:1.22.3 AS test
WORKDIR /cart
COPY ./cart .

FROM golang:1.22.3 AS builder
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GOARCH=amd64
WORKDIR /cart
COPY ./cart .
RUN go mod download
WORKDIR /cart/cmd
RUN go build -o /buildapp/cart
WORKDIR /cart

FROM alpine:latest AS production
ARG CARTAPP_TOPORT
RUN apk add --no-cache ca-certificates
COPY --from=builder /buildapp/cart /buildapp/cart
COPY ./cart/config ./buildapp/config
EXPOSE ${CARTAPP_TOPORT}
WORKDIR /buildapp
CMD ./cart