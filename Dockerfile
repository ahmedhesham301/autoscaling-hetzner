FROM golang:1.25.7-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build .

FROM alpine

WORKDIR /app

COPY --from=build /app/autoscaling-hetzner .

CMD [ "./autoscaling-hetzner" ]