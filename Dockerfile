FROM golang:1.25.7-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build .

FROM alpine:3

WORKDIR /app

COPY --from=build /app/autoscaling-hetzner .

CMD [ "./autoscaling-hetzner" ]