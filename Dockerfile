FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY tailscalehelper/. ./tailscalehelper/
COPY server/. ./server/

RUN CGO_ENABLED=0 GOOS=linux go build -o /command-center

FROM golang:alpine
COPY --from=builder /command-center /command-center
CMD [ "/command-center" ]

EXPOSE 3456
VOLUME [ "/data" ]