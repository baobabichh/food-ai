#!/bin/bash

stop_service() {
  local script_name=$1

  # Use pgrep to find the PID of the running service
  local pid=$(pgrep -f "./$script_name")

  # Check if the PID was found and kill the process
  if [ -n "$pid" ]; then
    echo "Stopping $script_name with PID: $pid"
    kill "$pid"
  else
    echo "No running instance of $script_name found."
  fi
}

current_dir=$(pwd)

cd "$current_dir/executable"
stop_service "food_recognition_exe.sh"

cd "$current_dir/executable"
stop_service "web_server_exe.sh"

cd "$current_dir/executable"
stop_service "telegram_bot_exe.sh"