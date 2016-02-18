cmake_minimum_required(VERSION 2.8.3)
project(rosgo)

find_package(catkin REQUIRED COMPONENTS rosgo_messages)

catkin_package(CFG_EXTRAS rosgo-extras.cmake)

catkin_add_env_hooks(99.rosgo SHELLS sh
                     DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/env-hooks)

include(${PROJECT_SOURCE_DIR}/cmake/rosgo-extras.cmake)

catkin_add_go_library(xmlrpc TARGET ${PROJECT_NAME}_xmlrpc)

catkin_add_go_library(
    ros
    TARGET ${PROJECT_NAME}
    DEPENDS ${PROJECT_NAME}_xmlrpc ${catkin_EXPORTED_TARGETS}
)

catkin_add_go_test(xmlrpc DEPENDS ${PROJECT_NAME}_xmlrpc)
catkin_add_go_test(ros DEPENDS ${PROJECT_NAME})

