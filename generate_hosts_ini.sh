#!/bin/bash


HOSTS_FILE_PATH="Ansible/hosts.ini"

master_ips=`multipass list | grep "arvan-challenge" | grep master | sort -k1 | awk '{ print $3 }'`
worker_ips=$(multipass list | grep "arvan-challenge" | grep worker | sort -k1 | awk '{ print $3 }')

rm -rf $HOSTS_FILE_PATH

echo "[masters]" | tee -a $HOSTS_FILE_PATH

i=1
for ip in "${master_ips[@]}"; do
	echo master-node-${i}		ansible_host=${ip}	ansible_user=ubuntu | tee -a $HOSTS_FILE_PATH
	i=$(( i + 1 ))
done

echo | tee -a $HOSTS_FILE_PATH
echo "[workers]" | tee -a $HOSTS_FILE_PATH

i=1
for ip in `echo "${worker_ips[@]}"`; do
	echo worker-node-${i}		ansible_host=${ip}	ansible_user=ubuntu | tee -a $HOSTS_FILE_PATH
	i=$(( i + 1 ))
done
