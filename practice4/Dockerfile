FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp ./cmd/api/main.go

FROM alpine:3.20
COPY --from=builder /app/myapp /myapp
COPY --from=builder /app/database/migrations /database/migrations
ENTRYPOINT ["/myapp"]
