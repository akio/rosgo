cmake_minimum_required(VERSION 2.8.3)
project(rosgo_core)

find_package(catkin REQUIRED COMPONENTS rosgo_build_tools rosgo_messages)

catkin_package()

catkin_add_go_library(xmlrpc TARGET ${PROJECT_NAME}_xmlrpc)

catkin_add_go_library(
    ros
    TARGET ${PROJECT_NAME}
    DEPENDS ${PROJECT_NAME}_xmlrpc ${catkin_EXPORTED_TARGETS}
)

catkin_add_go_test(xmlrpc DEPENDS ${PROJECT_NAME}_xmlrpc)
catkin_add_go_test(ros DEPENDS ${PROJECT_NAME})

