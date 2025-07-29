FROM golang:1.23

WORKDIR /opt/banners

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . "/opt/banners/"
