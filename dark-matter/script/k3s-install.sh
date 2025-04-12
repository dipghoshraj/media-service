#!/bin/bash
set -e

echo "Installing k3s on master..."

# Install k3s without Traefik or extra stuff
curl -sfL https://get.k3s.io | sh -s - --disable traefik --disable servicelb

echo "✅ k3s master installed"
echo "Token (save this for worker nodes):"
sudo cat /var/lib/rancher/k3s/server/node-token
