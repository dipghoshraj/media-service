#!/bin/bash
# This script is used to build a Kubernetes cluster using kubspray.
# It will install the necessary packages, configure the system, and initialize the cluster.
# Usage: ./k8s-cluster-build.sh <master-ip> <worker-ip>

# Check if the script is run as root
# if [ "$EUID" -ne 0 ]; then
#     echo "Please run as root"
#     exit 1
# fi

echo "📦 Installing requirements..."
sudo apt update
sudo apt install -y python3-pip git sshpass
pip3 install --upgrade pip
pip3 install ansible


echo "📁 Cloning Kubespray..."
git clone https://github.com/kubernetes-sigs/kubespray.git
cd kubespray
sudo pip3 install -r requirements.txt


echo "🔧 Configuring Ansible..."
cp -rfp inventory/sample inventory/mycluster
cp ../inventory/mycluster/inventory.ini inventory/mycluster/inventory.ini

# echo "🔑 Generating SSH keys..."
# ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -N ""
# ssh-copy-id -i ~/.ssh/id_rsa.pub root@${1}
# ssh-copy-id -i ~/.ssh/id_rsa.pub root@${2}
# echo "🔑 Copying SSH keys to all nodes..."

# for host in ${1} ${2}; do
#     sshpass -p 'password' ssh-copy-id -i ~/.ssh/id_rsa.pub root@${host}
# done
# echo "🔑 SSH keys copied to all nodes."

echo "🚀 Running the playbook..."
ansible-playbook -i inventory/${CLUSTER_NAME}/inventory.ini \
  --private-key=${SSH_PRIVATE_KEY} \
  -u ${SSH_USER} \
  cluster.yml -b -v

echo "✅ Kubernetes cluster is set up!"