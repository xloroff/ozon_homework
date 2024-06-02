FROM golang:1.22 AS e2etest
ARG GOAPP_TOPORT
COPY . .
COPY ./config /config
COPY ./log /log
WORKDIR /cart
RUN go mod download
RUN go install github.com/cespare/reflex@latest
EXPOSE ${GOAPP_TOPORT?goapp_toport_notset}
CMD reflex -r '(\.go$|go\.mod)' make run-all --start-service