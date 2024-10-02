#!/bin/bash

INTERFACE=eth0

## Get the IP address of the host
HOST_IP=$(ip address show dev $INTERFACE  | grep -oP '(?:\b\.?(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){4}/\d\d' | head -1)

echo "Host IP: $HOST_IP"

## Create the proxmox bridge
#pvesh create /nodes/pve/network --iface vmbr1 --type bridge --bridge_ports $INTERFACE --cidr "$HOST_IP"
#
## Apply the changes
#pvesh set /nodes/pve/network
