#!/bin/bash

OTHER_PID=0

cleanup() {
  echo "SIGTERM received. Performing cleanup..."
  kill $OTHER_PID
  exit 0
}
trap cleanup SIGTERM

cd ../services/python/
python3 food_recognition.py &


OTHER_PID=$!
wait $OTHER_PID

