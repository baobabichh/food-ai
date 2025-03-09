#!/bin/bash

OTHER_PID=0

cleanup() {
  echo "SIGTERM received. Performing cleanup..."
  kill $OTHER_PID
  kill 0
  exit 0
}
trap cleanup SIGTERM SIGINT

cd ../web_server/build/
./web_server &


OTHER_PID=$!
wait $OTHER_PID