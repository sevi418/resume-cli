FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/resume-cli .

FROM alpine:3.22

RUN apk add --no-cache poppler-utils ca-certificates
COPY --from=builder /out/resume-cli /usr/local/bin/resume-cli

ENTRYPOINT ["resume-cli"]
