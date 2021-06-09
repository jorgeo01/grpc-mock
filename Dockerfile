# Stage 1 - the build process
FROM golang:1.16 AS build-env

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" ./cmd/grpc-mock

# Stage 2 - the production environment
FROM alpine

COPY --from=build-env /app/grpc-mock /bin/

EXPOSE 22222

CMD ["grpc-mock"]
