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

sudo apt-get -y install golang-go 
sudo apt update 
sudo apt upgrade

go get github.com/PaulSonOfLars/gotgbot/v2

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

cd $current_dir
cd services/go
go mod init telegram_bot
go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5

cd $current_dir