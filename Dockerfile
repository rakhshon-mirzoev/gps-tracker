FROM golang:alpine AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.io,https://goproxy.cn,direct

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o core ./cmd/app

FROM alpine

WORKDIR /cmd

COPY --from=builder /app/core /cmd/

RUN chmod +x core

EXPOSE 8080

CMD ["./core"]

#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