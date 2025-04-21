FROM golang:1.24.0-alpine AS builder

COPY . /github.com/totorialman/go-task-avito
WORKDIR /github.com/totorialman/go-task-avito

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main/main.go
RUN go clean --modcache

FROM scratch AS runner

WORKDIR /build_v1/

COPY --from=builder /github.com/totorialman/go-task-avito/.bin .

COPY --from=builder /github.com/totorialman/go-task-avito/internal/middleware/acl/model.conf ./internal/middleware/acl/model.conf
COPY --from=builder /github.com/totorialman/go-task-avito/internal/middleware/acl/policy.csv ./internal/middleware/acl/policy.csv

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8080

ENTRYPOINT ["./.bin"]