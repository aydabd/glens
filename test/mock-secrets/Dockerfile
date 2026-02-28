FROM golang:1.25-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /mock-secrets .

FROM alpine:3.19
COPY --from=build /mock-secrets /mock-secrets
EXPOSE 8088
ENTRYPOINT ["/mock-secrets"]
