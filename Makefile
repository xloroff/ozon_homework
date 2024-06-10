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

run:
	docker-compose up -d

e2e-test: run-all run-e2e test stop

stop:
	docker-compose down

test:
	docker exec appe2e-production bash -c "go test -vet=off -v -race -count 1 -tags e2e ./internal/pkg/e2e/..."