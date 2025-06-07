#!/bin/bash
set -e

source ../env.sh

if [ -z "$1" ]; then
  echo "Usage: $0 <NODE_TOKEN>"
  exit 1
fi

NODE_TOKEN=$1

echo "Installing k3s agent on worker..."
sudo curl -sfL https://get.k3s.io | K3S_URL=https://${MASTER_IP}:6443 K3S_TOKEN=${NODE_TOKEN} sh -
echo "✅ k3s agent installed and connected to master"
