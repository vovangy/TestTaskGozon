FROM golang:1.21.0-alpine AS builder

COPY go.mod go.sum /myHabr/
WORKDIR /myHabr/

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./main ./cmd/main/main.go

FROM alpine:latest

WORKDIR /docker-myHabr/

COPY --from=builder /myHabr/config config/
COPY --from=builder /myHabr/main .

EXPOSE 80 443

ENTRYPOINT ["./main"]