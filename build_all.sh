#!/bin/bash
current_dir=$(pwd)

cd $current_dir
cd services/cpp
./build_all.sh

cd $current_dir