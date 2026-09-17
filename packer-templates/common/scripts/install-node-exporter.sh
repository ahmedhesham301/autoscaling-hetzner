#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
set -euxo pipefail

apt-get install prometheus-node-exporter -y
systemctl enable prometheus-node-exporter