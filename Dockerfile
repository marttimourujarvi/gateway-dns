FROM golang:1.27.1 AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/gateway-dns/main.go

FROM gcr.io/distroless/base:latest

COPY --chmod=755 --from=build /build/main /main

ENTRYPOINT ["/main"]

EXPOSE 53

