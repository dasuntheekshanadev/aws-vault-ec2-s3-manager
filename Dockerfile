FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ec2s3-cli .

FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /app/ec2s3-cli .

RUN apk add --no-cache ca-certificates

ENTRYPOINT ["./ec2s3-cli"]

