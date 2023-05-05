FROM golang:alpine AS builder

WORKDIR /project/ssh-honeypot

COPY go.* ./

RUN go mod download

COPY . .
RUN go build -o /project/ssh-honeypot/build/main .

FROM alpine:latest
COPY --from=builder /project/ssh-honeypot/build/main /app/build/main

EXPOSE 2222
ENTRYPOINT [ "/app/build/main" ]
