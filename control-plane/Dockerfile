FROM golang:1.27-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o control-plane-bin ./control-plane

FROM alpine:3

WORKDIR /app

COPY --from=build /app/control-plane-bin .

EXPOSE 8085

CMD [ "./control-plane-bin" ]