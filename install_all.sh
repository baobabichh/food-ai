#!/bin/bash
current_dir=$(pwd)

sudo apt update
sudo apt upgrade

sudo apt-get install -y cmake
sudo apt update 
sudo apt upgrade

sudo apt install -y clang
sudo apt update 
sudo apt upgrade

sudo apt install -y python3
sudo apt update 
sudo apt upgrade

sudo apt install -y python3-pip
sudo apt update 
sudo apt upgrade

sudo apt install -y mysql-server
sudo apt update 
sudo apt upgrade

sudo service mysql restart
sudo systemctl restart mysql.service

python3 -m pip config set global.break-system-packages true
pip install openai
pip install mysql-connector-python

libs_path="libs"
full_libs_path="$current_dir/$sub_path"
cd $full_libs_path
./install_drogon.sh

cd $full_libs_path
./install_nlohmanjson.sh