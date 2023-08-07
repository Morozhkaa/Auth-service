FROM golang:1.19-buster as builder

WORKDIR /auth
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o application

FROM alpine:3.15.4
WORKDIR /auth
COPY --from=builder /auth/application /auth/application
COPY *.yaml ./
CMD ["/auth/application"]

