#!/bin/bash

# Uninstall script for Table View

#----------------------------------------------------------------------

# Check user is root otherwise exit script

if [ "$EUID" -ne 0 ]
then
  printf "\nPlease run as root\n\n";
  exit;
fi;

cd /root;

#----------------------------------------------------------------------

# Stop Table View automatically starting on boot

systemctl stop tableview.service;
systemctl disable tableview.service;

# Remove Table View unit file and reload systemd deamon

rm /usr/lib/systemd/system/tableview.service;
systemctl daemon-reload;

#----------------------------------------------------------------------

# Remove Table View binary

rm /usr/bin/tableview;

# Remove all other directores and files used by Table View

rm -r /etc/tableview;

# Remove Table View source code in root home directory

rm -r /root/go/src/tableview;

# Remove the user and group tableview from the system

userdel tableview;
