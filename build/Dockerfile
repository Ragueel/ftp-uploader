FROM golang:1.21.1-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /var/bin/ftp-uploader ./cmd/ftp-uploader/main.go
RUN chmod +x /var/bin/ftp-uploader
ENV PATH "$PATH:/var/bin"

CMD ["ftp-uploader"]