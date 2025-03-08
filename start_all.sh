
#!/bin/bash
current_dir=$(pwd)

cd $current_dir
nohup ./food_recognition_exe.sh

cd $current_dir
nohup ./web_server_exe.sh