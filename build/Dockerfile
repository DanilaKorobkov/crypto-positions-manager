FROM golang:1.23-bullseye AS builder

WORKDIR /bin

COPY .. .

# Описание флагов:
# https://pkg.go.dev/cmd/link
# -ldflags '-s -w' - негласный стандарт для продуктовых бинарников, позволяет сократить размер
# бинарников примерно на 30%, удаляя всю отладочную информацию из них.
RUN CGO_ENABLED=0 go build -ldflags '-s -w' -o ./ ./cmd/watcher

FROM alpine:3.14 AS runtime
WORKDIR /app


# Создаем системного пользователя
RUN adduser -S app
USER app

COPY --from=builder /bin/dsadas ./

CMD ["./watcher"]
