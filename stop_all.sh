#!/bin/bash
current_dir=$(pwd)

cd $current_dir
kill ./food_recognition_exe.sh

cd $current_dir
kill ./web_server_exe.sh