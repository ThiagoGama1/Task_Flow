FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o taskflow ./cmd/server

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=America/Sao_Paulo

WORKDIR /app

COPY --from=builder /app/taskflow .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 3000

CMD ["./taskflow"]
