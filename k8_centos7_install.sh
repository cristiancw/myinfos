#!/bin/bash

########
# Log the message.
#
# Input:
# $1 - message
#
# Output:
function logit {
	echo -e "$(date +'%Y-%m-%d %H:%M:%S.%N') - ${USER} - INFO - ${*}"
}

########
# Disable swap / disable selinux / disable firewalld.
#
# Input:
#
# Output:
function prepareEnviroment() {
    logit "Preparing the enviroment..."
    
    logit "....Disabled swap"
    swapoff -a
    sed -ie 's:\(.*\)\s\(swap\)\s\s*\(\w*\)\s*\(\w*\)\s*\(.*\):# \1 \2 \3 \4 \5:' /etc/fstab
    
    logit "....Disabled selinux"
    setenforce 0
    sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
    
    logit "....Disabled firewall"
    systemctl disable firewalld && systemctl stop firewalld

    logit "....Disabled ipv6 for kubernetes in file: /etc/sysctl.d/k8s.conf"
    cat <<EOF > /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
}

########
# Install the docker and enable the service.
#
# Input:
#
# Output:
function installDocker() {
    logit "Installing the docker..."
    yum install -y docker
    systemctl enable docker && systemctl start docker 
}

########
# Install the kubelet, kubeadm, kubectl, kubernetes-cni and bash-completion.
#
# Input:
#
# Output:
function installK8() {
    logit "Installing the kubernetes..."
    cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
    yum install -y kubelet kubeadm kubectl kubernetes-cni bash-completion
    systemctl enable kubelet && systemctl start kubelet
    echo "source <(kubeadm completion bash)" >> ~/.bashrc
    echo "source <(kubectl completion bash)" >> ~/.bashrc    
}

##############
#### MAIN ####
##############

prepareEnviroment
installDocker
installK8

exit 0
