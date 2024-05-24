FROM golang:1.21.0-alpine AS builder

COPY . /myHabr/
WORKDIR /myHabr/

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./main ./cmd/user/main.go

FROM alpine:latest

WORKDIR /docker-myHabr-user/

COPY --from=builder /myHabr/config config/
COPY --from=builder /myHabr/main .

EXPOSE 80 443

ENTRYPOINT ["./main"]