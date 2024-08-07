FROM golang:1.22-alpine AS build

WORKDIR /build
COPY . .
RUN go build -o server ./application/grpc

FROM alpine:latest AS run
EXPOSE 8080
WORKDIR /app
COPY --from=build /build/server .
ENTRYPOINT ["/app/server"]
