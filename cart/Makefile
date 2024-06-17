# Путь установки бинарей плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin

# Путь до бинарей при локальном запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc


run:
	go run ./cmd

run-all:
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
	@go test -vet=off ./... -coverprofile ../cover/cover_cart.out.tmp > /dev/null
	@grep -vE "_mock.go" ../cover/cover_cart.out.tmp > ../cover/cover_cart.out
	@rm -f ../cover/cover_cart.out.tmp
	@go tool cover -func ../cover/cover_cart.out | fgrep total | awk '{print "Покрытие тестами всего проекта: " $$3}'
	@rm -f ../cover/cover_cart.out

cover-html:
	@go test -vet=off -coverprofile ../cover/cover_cart.out.tmp ./... > /dev/null
	@grep -vE "_mock.go" ../cover/cover_cart.out.tmp > ../cover/cover_cart.out
	@rm -f ../cover/cover_cart.out.tmp
	@go tool cover -html ../cover/cover_cart.out -o ../cover/cover_cart.html
	@rm -f ../cover/cover_cart.out

cover-pkg:
	@go test -vet=off ./... -coverpkg='./internal/api/...','./internal/repository/...','./internal/service/...' \
		-coverprofile ../cover/cover_cart.out.tmp > /dev/null
	@grep -vE "_mock.go" ../cover/cover_cart.out.tmp > ../cover/cover_cart.out
	@rm -f ../cover/cover_cart.out.tmp
	@go tool cover -func ../cover/cover_cart.out | fgrep total | awk '{print "Покрытие тестами уровня пакетов (api, service, repository): " $$3}'
	@rm -f ../cover/cover_cart.out

PHONY: .lint .test .bench

e2e-test-run:
	go test -vet=off -v -tags e2e ./...

PHONY: genproto
genproto: .proto-generate
	go mod tidy

PHONY: .proto-generate
.proto-generate: .bin-proto .vendor-proto .order-api-generate .stock-api-generate  .merge-swagger

# Вендоринг внешних proto файлов
.vendor-proto: .vendor-rm  vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options
	go mod tidy

# Чистка файлов вендоринга
.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto

# Устанавливаем необходимые бинарники для компиляции
# Для начала установить protoc через brew install protobuf или по-другому под вашу ОС
.PHONY: .bin-proto
.bin-proto:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2 && \
    GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/validate
		mv vendor-proto/tmp/validate vendor-proto/
		rm -rf vendor-proto/tmp

# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem

# Генерация кода через protoc
ORDER_API_PROTO_PATH:=api/order/v1
PHONY: .order-api-generate
.order-api-generate:
	rm -rf internal/pb/${ORDER_API_PROTO_PATH}
	mkdir -p internal/pb/${ORDER_API_PROTO_PATH}
	protoc \
	-I ../loms/${ORDER_API_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out internal/pb/${ORDER_API_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out internal/pb/${ORDER_API_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:internal/pb/${ORDER_API_PROTO_PATH}" \
	../loms/${ORDER_API_PROTO_PATH}/*.proto

STOCK_API_PROTO_PATH:=api/stock/v1
PHONY: .stock-api-generate
.stock-api-generate:
	rm -rf internal/pb/${STOCK_API_PROTO_PATH}
	mkdir -p internal/pb/${STOCK_API_PROTO_PATH}
	protoc \
	-I ../loms/${STOCK_API_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out internal/pb/${STOCK_API_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out internal/pb/${STOCK_API_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:internal/pb/${STOCK_API_PROTO_PATH}" \
	../loms/${STOCK_API_PROTO_PATH}/*.proto