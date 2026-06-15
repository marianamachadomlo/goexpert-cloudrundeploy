FROM golang:1.26-alpine AS build

ARG SERVICE=server

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux GOARCH=amd64 go build -o server ./cmd/${SERVICE}

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/server .
ENTRYPOINT ["./server"]
