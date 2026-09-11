FROM golang:1.27-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o db-control-plane-bin ./db-control-plane

FROM ubuntu:resolute

RUN apt-get update && apt-get install wget gnupg -y && \
    wget -O - https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(grep -oP '(?<=UBUNTU_CODENAME=).*' /etc/os-release || lsb_release -cs) main" | tee /etc/apt/sources.list.d/hashicorp.list && \
    apt-get update && apt-get install packer

RUN packer plugins install github.com/hetznercloud/hcloud

WORKDIR /app

COPY --from=build /app/db-control-plane-bin .

EXPOSE 8090

CMD [ "./db-control-plane-bin" ]