FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY pkg/logging/ pkg/logging/
COPY cmd/api/ cmd/api/
WORKDIR /src/cmd/api
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /api .

FROM gcr.io/distroless/static
COPY --from=builder /api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
