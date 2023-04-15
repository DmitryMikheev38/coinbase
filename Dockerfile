FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/app/main.go

FROM alpine:latest
ENV APP_PORT=8080
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main ./
EXPOSE ${APP_PORT}

ENTRYPOINT ["./main"]