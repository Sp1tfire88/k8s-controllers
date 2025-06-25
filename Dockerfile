# Stage 1: Build with Go >= 1.23.0
FROM golang:1.24 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o controller main.go

# Stage 2: Distroless
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /app/controller .

USER nonroot:nonroot
ENTRYPOINT ["/controller"]
