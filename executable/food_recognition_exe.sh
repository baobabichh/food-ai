#!/bin/bash

OTHER_PID=0

cleanup() {
  echo "SIGTERM received. Performing cleanup..."
  kill $OTHER_PID
  kill 0
  exit 0
}
trap cleanup SIGTERM SIGINT

cd ../services/python/
python3 food_recognition.py &


OTHER_PID=$!
wait $OTHER_PID

