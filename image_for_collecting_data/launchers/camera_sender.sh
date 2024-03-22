#!/bin/bash

source /environment.sh

# initialize launch file
dt-launchfile-init

# launch publisher
rosrun camera_sender camera_sender_node.py

# wait for app to end
dt-launchfile-join