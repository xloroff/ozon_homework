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

build:
	go build -o ./bin/loms ./bin/loms

lint:
	golangci-lint run --config=./golangci.yaml

pprof:
	go tool pprof -http :9999 http://localhost:8084/debug/pprof/profile\?seconds=5

cover:
	@go test -vet=off ./... -coverprofile ../cover/cover_loms.out.tmp > /dev/null
	@grep -vE "_mock.go" ../cover/cover_loms.out.tmp > ../cover/cover_loms.out
	@rm -f ../cover/cover_loms.out.tmp
	@go tool cover -func ../cover/cover_loms.out | fgrep total | awk '{print "Покрытие тестами всего проекта: " $$3}'
	@rm -f ../cover/cover_loms.out

cover-html:
	@go test -vet=off -coverprofile ../cover/cover_loms.out.tmp ./... > /dev/null
	@grep -vE "_mock.go" ../cover/cover_loms.out.tmp > ../cover/cover_loms.out
	@rm -f ../cover/cover_loms.out.tmp
	@go tool cover -html ../cover/cover_loms.out -o ../cover/cover_loms.html
	@rm -f ../cover/cover_loms.out

cover-pkg:
	@go test -vet=off ./... -coverpkg='./internal/service/...' \
		-coverprofile ../cover/cover_loms.out.tmp > /dev/null
	@grep -vE "_mock.go" ../cover/cover_loms.out.tmp > ../cover/cover_loms.out
	@rm -f ../cover/cover_loms.out.tmp
	@go tool cover -func ../cover/cover_loms.out | fgrep total | awk '{print "Покрытие тестами уровня пакетов (service): " $$3}'
	@rm -f ../cover/cover_loms.out

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
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0 && \
    GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/go-swagger/go-swagger/cmd/swagger@v0.31.0
	GOBIN=$(LOCAL_BIN) go install github.com/ntauth/go-swagger-merger@v0.1.0

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
	rm -rf pkg/${ORDER_API_PROTO_PATH}
	mkdir -p pkg/${ORDER_API_PROTO_PATH}
	protoc \
	-I ${ORDER_API_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out pkg/${ORDER_API_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out pkg/${ORDER_API_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:pkg/${ORDER_API_PROTO_PATH}" \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
	--grpc-gateway_out pkg/${ORDER_API_PROTO_PATH} \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
	--openapiv2_out pkg/${ORDER_API_PROTO_PATH} \
	--openapiv2_opt logtostderr=true \
	${ORDER_API_PROTO_PATH}/*.proto

STOCK_API_PROTO_PATH:=api/stock/v1
PHONY: .stock-api-generate
.stock-api-generate:
	rm -rf pkg/${STOCK_API_PROTO_PATH}
	mkdir -p pkg/${STOCK_API_PROTO_PATH}
	protoc \
	-I ${STOCK_API_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out pkg/${STOCK_API_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out pkg/${STOCK_API_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:pkg/${STOCK_API_PROTO_PATH}" \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
	--grpc-gateway_out pkg/${STOCK_API_PROTO_PATH} \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
	--openapiv2_out pkg/${STOCK_API_PROTO_PATH} \
	--openapiv2_opt logtostderr=true \
	${STOCK_API_PROTO_PATH}/*.proto

# Соединяем оба описания прикладов в один файл
PHONY: .merge-swagger
.merge-swagger:
	rm -rf pkg/swagger
	mkdir -p pkg/swagger
	$(LOCAL_BIN)/go-swagger-merger \
	-o pkg/docs/swagger.json \
	pkg/${ORDER_API_PROTO_PATH}/order.swagger.json \
	pkg/${STOCK_API_PROTO_PATH}/stock.swagger.json

PHONY: gensqlc
gensqlc:
	sqlc generate