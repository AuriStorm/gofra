FROM golang:1.23.5

WORKDIR /app

COPY . ./

CMD ["go", "run", "./cmd/app/", "-port=8090", "-default-timeout-sec=5", "-max-queues=2", "-queue-size=5"]
