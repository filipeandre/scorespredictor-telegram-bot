#!/usr/bin/env bash

sudo yum update && sudo yum install wget nano git
wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go*
sudo timedatectl set-timezone Europe/Rome

#-------------------------
# /etc/profile
#-------------------------
sudo sh -c 'echo "export PATH=/usr/local/go/bin:\$PATH" > /etc/profile.d/go.sh'
source /etc/profile