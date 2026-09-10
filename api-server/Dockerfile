FROM golang:1.27-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build .

FROM alpine:3

WORKDIR /app

COPY --from=build /app/api-server .

EXPOSE 8080

CMD [ "./api-server" ]