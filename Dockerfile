FROM golang:1.26-alpine AS build
WORKDIR /app
# replace ../gateway-be из /app указывает на /gateway-be; нужен только модуль и сгенерированный код
COPY --from=gateway-be go.mod go.sum /gateway-be/
COPY --from=gateway-be gen /gateway-be/gen
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /users-be ./cmd/server

FROM alpine:latest
COPY --from=build /users-be /users-be
EXPOSE 50051
ENTRYPOINT ["/users-be"]
