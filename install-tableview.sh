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
  printf "Please run commands: \"cd /root; git clone https://github.com/ellwould/table-view\"\n";
  printf "and run install script again\n\n";
  exit;
fi;

#----------------------------------------------------------------------

# Copy unit files and reload systemd deamon

cp /root/table-view/systemd/tableview.service /usr/lib/systemd/system/;
systemctl daemon-reload;

#----------------------------------------------------------------------

# Install wget

apt update;
apt install wget;

#----------------------------------------------------------------------

# Remove any previous version of Go, download and install Go 1.24.5

wget -P https://go.dev/dl/go1.24.5.linux-amd64.tar.gz;
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz;
export PATH=$PATH:/usr/local/go/bin;

#----------------------------------------------------------------------

# Create HTML/CSS directory and copy HTML/CSS start and end file

mkdir -p /etc/tableview/html-css;
cp /root/table-view/html-css/* /etc/tableview/html-css/;

# Copy /root/table-view/env/tableview.env into /etc/tableview

cp /root/table-view/env/tableview.env /etc/tableview/tableview.env;

# Create Go directories in root home directory

mkdir -p /root/go/{bin,pkg,src/tableview};

# Copy Go source code

cp /root/table-view/go/tableview.go /root/go/src/tableview/tableview.go;

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
chmod 050 /etc/tableview;
chmod 040 /etc/tableview/tableview.env;
chmod 050 /etc/tableview/html-css;
chmod 040 /etc/tableview/html-css/*;

# Enable tableview on boot

systemctl enable tableview;

#----------------------------------------------------------------------

printf "\nUpdate database details in /etc/tableview/tableview.env\n";
printf "\nThen to start Table View run: systemctl start tableview\n";
