build-all:
	docker-compose up --force-recreate --build

run-all:
	docker-compose up -d --force-recreate

run-all-dev:
	docker-compose --env-file .env --env-file .env.dev up --force-recreate

run-e2e:
	docker-compose up -d e2e

run-e2e-dev:
	docker-compose --env-file .env --env-file .env.dev up -d e2e

run-all-test:
	docker-compose --env-file .env.test up -d loms

run:
	docker-compose up -d

e2e-test: run-all run-e2e test stop clear-volume

i-test: run-all-test integration-test stop clear-volume

clear-volume:
	docker volume rm $$(docker volume ls -qf dangling=true) || true

integration-test:
	docker exec apploms-test bash -c "go test -vet=off -v -race -count 1 -tags integration ./internal/pkg/integration/..."

stop:
	docker-compose down --remove-orphans || true

test:
	docker exec appe2e-production bash -c "go test -vet=off -v -race -count 1 -tags e2e ./internal/pkg/e2e/..."

migration-up:
	docker exec migration-production bash -c "go run ./cmd/main.go up"

migration-up-dev:
	docker exec migration-development bash -c "go run ./cmd/main.go up"
