# Dockerfile

# build stage
FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o controller main.go

# build stage (distroless)
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /app/controller .

USER nonroot:nonroot

EXPOSE 9000

ENTRYPOINT ["/controller", "server"]
