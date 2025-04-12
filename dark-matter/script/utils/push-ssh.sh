#!/bin/bash

# This script is used to push SSH keys to all nodes in the cluster.
# It will generate SSH keys if they do not exist and copy them to all nodes.
# Usage: ./push-ssh.sh <master-ip> <worker-ip>

set -e

SSH_USER="ubuntu"
PUB_KEY="$HOME/.ssh/id_rsa.pub"

if [ ! -f "$PUB_KEY" ]; then
  echo "❌ SSH public key not found."
  exit 1
fi

echo "🔑 Using public key at: $PUB_KEY"
read -s -p "Enter SSH password for initial access: " SSH_PASS
echo


while read -r IP; do
  echo "🔐 Pushing key to $IP"
  sshpass -p "$SSH_PASS" ssh -o StrictHostKeyChecking=no ${SSH_USER}@$IP \
    "mkdir -p ~/.ssh && echo '$(cat $PUB_KEY)' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys"
done