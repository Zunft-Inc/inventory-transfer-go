FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 go build -o /inventory-transfer-service .

FROM alpine:3.22
COPY --from=build /inventory-transfer-service /usr/local/bin/inventory-transfer-service
USER nobody
EXPOSE 8080
CMD ["inventory-transfer-service"]
