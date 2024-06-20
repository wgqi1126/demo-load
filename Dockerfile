FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /tmp/app

FROM ubuntu:24.04

WORKDIR /app

COPY --from=builder /tmp/app /app/app

EXPOSE 80

CMD ["/app/app"]
