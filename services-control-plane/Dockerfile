FROM golang:1.27-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o services-control-plane-bin ./services-control-plane

FROM alpine:3

WORKDIR /app

COPY --from=build /app/services-control-plane-bin .

EXPOSE 8090

CMD [ "./services-control-plane-bin" ]