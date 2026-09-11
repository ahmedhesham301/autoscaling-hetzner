FROM golang:1.27-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download


RUN go build -o api-server-bin ./api-server

FROM alpine:3

WORKDIR /app

COPY --from=build /app/api-server-bin .

EXPOSE 8080

CMD [ "./api-server-bin" ]