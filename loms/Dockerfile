FROM golang:1.22.3 AS development
WORKDIR /loms
ARG LOMSAPP_TOPORT
ARG LOMSHTTP_TOPORT
COPY ./loms .
RUN go mod download
RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      protobuf-compiler
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install github.com/gojuno/minimock/v3/cmd/minimock@latest
RUN go install github.com/cespare/reflex@latest
EXPOSE ${LOMSAPP_TOPORT?lomsapp_toport_notset} ${LOMSHTTP_TOPORT?lomshttp_toport_notset}
CMD reflex -r '(\.go$|go\.mod)' make run-all --start-service


FROM golang:1.22.3 AS builder
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GOARCH=amd64
WORKDIR /loms
COPY ./loms .
RUN go mod download
WORKDIR /loms/cmd
RUN go build -o /buildapp/loms
WORKDIR /loms

FROM alpine:latest AS production
ARG LOMSAPP_TOPORT
ARG LOMSHTTP_TOPORT
RUN apk add --no-cache ca-certificates
COPY --from=builder /buildapp/loms /buildapp/loms
COPY ./loms/config ./buildapp/config
COPY ./loms/stock-data.json ./buildapp
EXPOSE ${LOMSAPP_TOPORT} ${LOMSHTTP_TOPORT}
WORKDIR /buildapp
CMD ./loms