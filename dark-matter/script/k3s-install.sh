#!/bin/bash
set -e

echo "Installing k3s on master..."

# Install k3s without Traefik or extra stuff
curl -sfL https://get.k3s.io | sh -s - --disable traefik --disable servicelb --node-external-ip=15.207.221.101 --tls-san=15.207.221.101 --tls-san=ip-172-31-2-61


echo "✅ k3s master installed"
echo "Token (save this for worker nodes):"
sudo cat /var/lib/rancher/k3s/server/node-token
