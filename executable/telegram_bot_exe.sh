#!/bin/bash

OTHER_PID=0

cleanup() {
  echo "SIGTERM received. Performing cleanup..."
  kill $OTHER_PID
  kill 0
  exit 0
}
trap cleanup SIGTERM SIGINT

cd ../services/go/
go run telegram_bot.go &


OTHER_PID=$!
wait $OTHER_PID

