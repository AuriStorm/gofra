# GOFRA

## Yet another go message queue sample service

### способы запуска

локально:

```make dev```

локально с хотрелоадом для удобства дебага (локально должно быть установлено https://github.com/mitranim/gow):

```make gow```

в докер контейнере:

```make docker-dev```

локально без make:

```go run ./cmd/app/ -port=8090 -default-timeout-sec=5 -max-queues=2 -queue-size=5```
