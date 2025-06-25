# Stage 1: Build
FROM golang:1.21 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o controller main.go

# Stage 2: Minimal Distroless
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /app/controller .

USER nonroot:nonroot

EXPOSE 8080
ENTRYPOINT ["/controller"]
