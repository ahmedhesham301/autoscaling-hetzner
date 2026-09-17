#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
set -euxo pipefail

apt-get update
apt-get full-upgrade -y
