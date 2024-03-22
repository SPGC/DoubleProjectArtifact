cmake_minimum_required(VERSION 2.8.3)
project(camera_sender)

find_package(catkin REQUIRED COMPONENTS
  rospy
  cv_bridge
)

catkin_package()