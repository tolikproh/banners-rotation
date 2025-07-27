FROM golang:1.23 AS build

ENV CGO_ENABLED=0
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.22

COPY --from=build /go/bin/goose /bin/goose

WORKDIR /var/migrations

COPY migrations/*.sql /var/migrations/

CMD ["/bin/goose", "-dir", ".", "up" ]