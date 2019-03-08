#!/usr/bin/env bash

sudo yum update && sudo yum install wget nano
wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go*
sudo timedatectl set-timezone Europe/Rome

#-------------------------
# /etc/profile
#-------------------------
sudo sh -c ' cat <<EOT >> /etc/profile
export SCORESPREDICTOR_HOME=/home/nicola/source/scorespredictor/
EOT
' && source /etc/profile

#------------------------------
# /etc/cron.d/scorespredictor
#------------------------------
sudo sh -c ' cat <<EOT > /etc/cron.d/scorespresdictor
SHELL=/bin/bash
MAILTO=nicola
CRON_TZ=Europe/Rome

0 10 * * * nicola source /etc/profile && cd /home/nicola/source/scorespredictor/ && go run *go
EOT
' && sudo systemctl restart crond.service