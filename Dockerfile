FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY tailscale/*.go ./tailscale/
COPY assets assets

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-nixos-menu

FROM golang:alpine
COPY --from=builder /go-nixos-menu /go-nixos-menu
CMD [ "/go-nixos-menu" ]

EXPOSE 3000
VOLUME [ "/data" ]