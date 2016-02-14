cmake_minimum_required(VERSION 2.8.3)
project(rosgo_core)

find_package(catkin REQUIRED COMPONENTS rosgo_build_tools gengo)

catkin_package()

rosgo_generate_messages(std_msgs rosgraph_msgs)

catkin_add_go_library(xmlrpc TARGET ${PROJECT_NAME}_xmlrpc)

catkin_add_go_library(
    ros
    TARGET ${PROJECT_NAME}
    DEPENDS ${PROJECT_NAME}_xmlrpc std_msgs_gengo rosgraph_msgs_gengo
)

catkin_add_go_test(xmlrpc DEPENDS ${PROJECT_NAME}_xmlrpc)
catkin_add_go_test(ros DEPENDS ${PROJECT_NAME})

