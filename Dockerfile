FROM golang:alpine AS builder

WORKDIR /project/ssh-honeypot

COPY go.* ./

RUN go mod download

COPY . .
RUN go build -o build/main cmd/ssh-honeypot/main.go

FROM alpine:latest

RUN apk update && apk add --no-cache openssh-keygen

COPY --from=builder /project/ssh-honeypot/build/main /app/build/main

COPY ${CMDS_FILE} ${CMDS_FILE}

EXPOSE 2222

ENTRYPOINT [ "/app/build/main" ]
