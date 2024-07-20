FROM golang:1.22.3 AS development
WORKDIR /notifier
COPY ./notifier .
RUN go mod download
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install github.com/cespare/reflex@latest
CMD reflex -r '(\.go$|go\.mod)' make run-all --start-service

FROM golang:1.22.3 AS e2e
WORKDIR /notifier
COPY ./notifier .


FROM golang:1.22.3 AS test
WORKDIR /notifier
COPY ./notifier .

FROM golang:1.22.3 AS builder
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GOARCH=amd64
WORKDIR /notifier
COPY ./notifier .
RUN go mod download
WORKDIR /notifier/cmd
RUN go build -o /buildapp/notifier
WORKDIR /notifier

FROM alpine:latest AS production
RUN apk add --no-cache ca-certificates
COPY --from=builder /buildapp/notifier /buildapp/notifier
COPY ./notifier/config ./buildapp/config
WORKDIR /buildapp
CMD ./notifier