#!/bin/bash
set -e

CLUSTER_NAME="mycluster"
INVENTORY_DIR="./inventory/${CLUSTER_NAME}"
WORKERS_INI="${INVENTORY_DIR}/workers.ini"
SSH_KEY="$HOME/.ssh/id_rsa"
SSH_USER="ubuntu"

echo "📂 Preparing worker inventory..."
cp ../${WORKERS_INI} kubespray/inventory/${CLUSTER_NAME}/inventory.ini

echo "🚀 Adding worker nodes to the cluster..."
cd kubespray
ansible-playbook -i inventory/${CLUSTER_NAME}/inventory.ini \
  --private-key=${SSH_KEY} \
  -u ${SSH_USER} \
  scale.yml -b -v

echo "✅ Workers successfully added to the cluster!"
