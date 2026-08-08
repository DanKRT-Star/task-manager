# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

# ---- Run stage ----
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 3000

CMD ["./server"]