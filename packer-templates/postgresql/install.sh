#!/bin/bash
set -eux
export DEBIAN_FRONTEND=noninteractive

# verify number or args
if [[ "$#" -ne 2 ]]; then
    echo "wrong number of args"
    exit 1
fi

if ! [[ "$1" =~ ^[0-9]+$ ]]; then
    echo "First argument must be a number got $1"
    exit 1
fi

if [[ "$2" != "true" && "$2" != "false" ]]; then
    echo "Second argument must be true or false not $1"
    exit 1
fi

APP_VERSION="$1"
NODE_EXPORTER="$2"
# update and upgrade
apt-get update
apt-get upgrade -y
apt-get autopurge -y


# Disable automatic service startup
echo -e '#!/bin/sh\nexit 101' > /usr/sbin/policy-rc.d
chmod +x /usr/sbin/policy-rc.d

# setup postgresql repo
apt-get install postgresql-common -y
/usr/share/postgresql-common/pgdg/apt.postgresql.org.sh -y
apt-get update


# install
apt-get install "postgresql-$APP_VERSION" etcd-server etcd-client python3-etcd3 python3-etcd patroni -y

if [[ $NODE_EXPORTER == "true" ]]; then
    apt-get install prometheus-node-exporter -y
    systemctl enable prometheus-node-exporter
fi

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
