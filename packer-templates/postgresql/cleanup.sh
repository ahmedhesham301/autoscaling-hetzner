#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
set -eux

# Enable automatic service startup
rm /usr/sbin/policy-rc.d

# Clean and reset could-init files
cloud-init clean --logs --machine-id --seed --configs all

rm -rf /run/cloud-init/*
rm -rf /var/lib/cloud/*

# Clean apt files

apt-get -y autopurge
apt-get -y clean

rm -rf /var/lib/apt/lists/*

# Clean logs
journalctl --flush
journalctl --rotate --vacuum-time=0

find /var/log -type f -exec truncate --size 0 {} \; # truncate system logs
find /var/log -type f -name '*.[1-9]' -delete # remove archived logs
find /var/log -type f -name '*.gz' -delete # remove compressed archived logs

rm -rf /var/log/postgresql/*

# Reset host ssh keys
rm -f /etc/ssh/ssh_host_*_key /etc/ssh/ssh_host_*_key.pub

# discard the now unused blocks from the disk
fstrim --all || true
sync
