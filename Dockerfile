FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./

RUN go build -o /app/cmd/server/main ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/cmd/server/main /main
EXPOSE 8080
CMD ["/main"]
