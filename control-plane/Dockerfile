FROM golang:1.27-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build .

FROM alpine:3

WORKDIR /app

COPY --from=build /app/control-plane .

EXPOSE 8085

CMD [ "./control-plane" ]