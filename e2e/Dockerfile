FROM golang:1.22.3 AS development
WORKDIR /e2e
COPY ./e2e .
RUN go mod download

FROM golang:1.22.3 AS production
WORKDIR /e2e
COPY ./e2e .
RUN go mod download