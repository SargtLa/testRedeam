build:
	docker-compose up -d --build
	go test -v ./... > last_test.log