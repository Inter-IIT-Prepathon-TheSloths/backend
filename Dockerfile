FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main ./cmd/main/main.go

FROM alpine:3.18

COPY --from=build /main /main

RUN ls -l /main

CMD ["/main"]
