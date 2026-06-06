FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /app/server .

EXPOSE 6019

CMD ["./server"]
