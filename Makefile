build-all:
	docker-compose up --force-recreate --build

run-all:
	docker-compose up --force-recreate

run:
	docker-compose up -d

stop:
	docker-compose down
	#docker rmi appcart-dev
