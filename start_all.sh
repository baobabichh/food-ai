#!/bin/bash

start_service() {
  local script_name=$1
  local log_file="${script_name%.*}.log"

  [ -f "$log_file" ] && rm "$log_file"
  nohup ./"$script_name" > "$log_file" 2>&1 &
  echo $!
}

current_dir=$(pwd)

cd "$current_dir/executable"
start_service "food_recognition_exe.sh" 

cd "$current_dir/executable"
start_service "web_server_exe.sh" 



