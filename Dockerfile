FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /pismo-api ./cmd/api

FROM alpine:3.22
WORKDIR /app
COPY --from=build /pismo-api /usr/local/bin/pismo-api
EXPOSE 8080
CMD ["pismo-api"]
