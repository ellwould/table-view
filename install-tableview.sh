#!/bin/bash

# Install Script for Table View

#----------------------------------------------------------------------

# Check user is root otherwise exit script

if [ "$EUID" -ne 0 ]
then
  printf "\nPlease run as root\n\n";
  exit;
fi;

cd /root;

#----------------------------------------------------------------------

# Check Table View has been cloned from GitHub

if [ ! -d "/root/table-view" ]
then
  printf "\nDirectory table-view does not exist in /root.\n";
  printf "Please run commands: \"cd /root; git clone https://github.com/Ellwould/table-view\"\n";
  printf "and run install script again\n\n";
  exit;
fi;

#----------------------------------------------------------------------

# Copy unit files and reload systemd deamon

cp /root/table-view/systemd/tableview.service /usr/lib/systemd/system/;
systemctl daemon-reload;

#----------------------------------------------------------------------

# Remove any previous version of Go, download and install Go 1.23.3

wget -P /root https://go.dev/dl/go1.23.3.linux-amd64.tar.gz;
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz;

#----------------------------------------------------------------------

# Create HTML/CSS directory and copy HTML/CSS start and end file

mkdir /usr/local/etc/tableview-resource;
cp /root/table-view/html-css/tableview-start.html /usr/local/etc/tableview-resource/;
cp /root/table-view/html-css/tableview-end.html /usr/local/etc/tableview-resource/;

# Copy /root/table-view/env/tableview.env into /usr/local/etc/tableview-resource

cp /root/table-view/env/tableview.env /usr/local/etc/tableview-resource/tableview.env;

# Create Go directories in root home directory

mkdir -p /root/go/{bin,pkg,src/tableview};

# Create tableviewresource Go directory

mkdir /usr/local/go/src/tableviewresource;

# Copy Go source code

cp /root/table-view/go/tableview.go /root/go/src/tableview/tableview.go;
cp /root/table-view/go/tableviewresource.go /usr/local/go/src/tableviewresource/tableviewresource.go;

# Create Go mod for tableview

export PATH=$PATH:/usr/local/go/bin;
cd /root/go/src/tableview;
go mod init root/go/src/tableview;
go mod tidy;

# Compile tableview.go

cd /root/go/src/tableview;
go build tableview.go;
cd /root;

# Create system user named tableview with no shell, no home directory and lock account

useradd -r -s /bin/false tableview;
usermod -L tableview;

# Change executables file permissions, owner, group and move executables

chown root:tableview /root/go/src/tableview/tableview;
chmod 050 /root/go/src/tableview/tableview;
mv /root/go/src/tableview/tableview /usr/local/bin/tableview;

# Change tableviewresource file permissions, owner and group

chown -R root:tableview /usr/local/etc/tableview-resource;
chmod 050 /usr/local/etc/tableview-resource;
chmod 040 /usr/local/etc/tableview-resource/*;

# Enable tableview on boot

systemctl enable tableview;

#----------------------------------------------------------------------

printf "\nUpdate database details in /usr/local/etc/tableview-resource/tableview.env\n";
printf "\nThen to start Table View run: systemctl start tableview\n";
