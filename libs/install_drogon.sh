sudo apt-get update -y
sudo apt-get upgrade -y

sudo apt install -y git gcc g++ cmake
sudo apt-get update -y
sudo apt-get upgrade -y

sudo apt install -y libjsoncpp-dev uuid-dev zlib1g-dev openssl libssl-dev
sudo apt-get update -y
sudo apt-get upgrade -y

sudo apt-get install -y libmysqlclient-dev
sudo apt-get update -y
sudo apt-get upgrade -y

sudo apt-get install -y libmariadb-dev-compat libmariadb-dev
sudo apt-get update -y
sudo apt-get upgrade -y

sudo apt-get install -y drogon
sudo apt-get update -y
sudo apt-get upgrade -y

rm -rf drogon/* drogon/.*
rmdir drogon
git clone https://github.com/drogonframework/drogon
cd drogon
git submodule update --init
mkdir build
cd build
cmake .. -DMYSQL=ON
make && sudo make install