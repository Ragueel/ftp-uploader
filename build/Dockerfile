FROM golang:1.21.1-alpine as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/ftp-uploader ./cmd/ftp-uploader/main.go
RUN chmod +x /bin/ftp-uploader

FROM alpine:latest as production

COPY --from=builder /bin/ftp-uploader /bin/ftp-uploader
RUN chmod +x /bin/ftp-uploader

CMD ["ftp-uploader"]