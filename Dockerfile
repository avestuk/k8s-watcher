FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

ENV USER=k8s-watcher
ENV UID=10001

RUN adduser \
	        --disabled-password \
	        --gecos "" \
	        --home "/nonexistent" \
	        --shell "/sbin/nologin" \
	        --no-create-home \
	        --uid "${UID}" \
	        "${USER}" && \
        apk update && apk add --no-cache git && \
        go mod download && \
        CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o k8s-watcher

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/k8s-watcher /k8s-watcher
USER k8s-watcher:k8s-watcher

ENTRYPOINT ["/k8s-watcher"]
