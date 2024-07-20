run:
	go run ./cmd

run-all: lint test cover
	go run ./cmd

test:
	go test -v -race -count 1 -vet=off ./...

bench:
	go test -v -bench=. -benchmem ./...

build:
	go build -o ./bin/cart ./bin/cart

lint:
	golangci-lint run --config=./golangci.yaml

cover:
	@go test -vet=off ./... -coverprofile ../cover/cover.out.tmp > /dev/null
	@grep -vE "_mock.go" ../cover/cover.out.tmp > ../cover/cover.out
	@rm -f ../cover/cover.out.tmp
	@go tool cover -func ../cover/cover.out | fgrep total | awk '{print "Покрытие тестами всего проекта: " $$3}'
	@rm -f ../cover/cover.out

cover-html:
	@go test -vet=off -coverprofile ../cover/cover.out.tmp ./... > /dev/null
	@grep -vE "_mock.go" ../cover/cover.out.tmp > ../cover/cover.out
	@rm -f ../cover/cover.out.tmp
	@go tool cover -html ../cover/cover.out -o ../cover/cover.html
	@rm -f ../cover/cover.out

cover-pkg:
	@go test -vet=off ./... -coverpkg='./internal/api/...','./internal/repository/...','./internal/service/...' \
		-coverprofile ../cover/cover.out.tmp > /dev/null
	@grep -vE "_mock.go" ../cover/cover.out.tmp > ../cover/cover.out
	@rm -f ../cover/cover.out.tmp
	@go tool cover -func ../cover/cover.out | fgrep total | awk '{print "Покрытие тестами уровня пакетов (api, service, repository): " $$3}'
	@rm -f ../cover/cover.out

PHONY: .lint .test .bench

e2e-test-run:
	go test -vet=off -v -tags e2e ./...
