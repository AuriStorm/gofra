
dev-build:
	go build -o bin/gofra ./cmd/app/

dev-run:
	./bin/gofra -port=8090 -default-timeout-sec=5 -max-queues=2 -queue-size=5

dev: dev-build dev-run

docker-dev-build:
	docker build -t gofra:dev .

# TODO pass port and defaulttimout params to container
docker-dev-run:
	docker run --rm -p 8090:8090 gofra:dev

docker-dev: docker-dev-build docker-dev-run

gow:
	gow run ./cmd/app/ -port=8090 -default-timeout-sec=5 -max-queues=2 -queue-size=5
