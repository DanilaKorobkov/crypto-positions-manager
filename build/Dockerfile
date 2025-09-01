FROM golang:1.23-bullseye AS builder

WORKDIR /bin

COPY .. .

RUN CGO_ENABLED=0 go build -ldflags '-s -w' -o ./watcher ./cmd/watcher && \
    CGO_ENABLED=0 go build -ldflags '-s -w' -o ./cli ./cmd/cli

FROM alpine:3.14 AS runtime

WORKDIR /app

# Создаем системного пользователя
RUN adduser -S app
USER app

COPY --from=builder --chown=app /bin/watcher /bin/cli ./

CMD ["./watcher"]
