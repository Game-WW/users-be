FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /users-be ./cmd/server

FROM alpine:latest
COPY --from=build /users-be /users-be
EXPOSE 50051
ENTRYPOINT ["/users-be"]
