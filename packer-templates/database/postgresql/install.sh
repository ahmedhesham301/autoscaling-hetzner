#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
set -euxo pipefail

# verify number or args
if [[ "$#" -ne 1 ]]; then
    echo "wrong number of args"
    exit 1
fi

if ! [[ "$1" =~ ^[0-9]+$ ]]; then
    echo "First argument must be a number got $1"
    exit 1
fi

APP_VERSION="$1"

# setup postgresql repo
apt-get install postgresql-common -y
/usr/share/postgresql-common/pgdg/apt.postgresql.org.sh -y
apt-get update


# install
apt-get install "postgresql-$APP_VERSION" etcd-server etcd-client python3-etcd3 python3-etcd patroni -y

# configure
pg_dropcluster $APP_VERSION main

systemctl edit --stdin patroni.service <<EOF
[Service]
User=postgres
Group=postgres
EOF

mkdir /data
chown -R postgres:postgres /data
systemctl daemon-reload

mkdir -p /etc/patroni
mv /tmp/patroni-config.yml /etc/patroni/config.yml

systemctl disable postgresql
systemctl enable etcd
systemctl enable patroni
