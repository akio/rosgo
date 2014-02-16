cmake_minimum_required(VERSION 2.8.3)
project(rosgo_core)

## Find catkin macros and libraries
## if COMPONENTS list like find_package(catkin REQUIRED COMPONENTS xyz)
## is used, also find other catkin packages
find_package(catkin REQUIRED COMPONENTS rosgo_build_tools gengo)

## System dependencies are found with CMake's conventions
# find_package(Boost REQUIRED COMPONENTS system)


## Uncomment this if the package has a setup.py. This macro ensures
## modules and global scripts declared therein get installed
## See http://ros.org/doc/groovy/api/catkin/html/user_guide/setup_dot_py.html
# catkin_python_setup()


#######################################
## Declare ROS messages and services ##
#######################################

## Generate messages in the 'msg' folder
# add_message_files(
#   FILES
#   Message1.msg
#   Message2.msg
# )

## Generate services in the 'srv' folder
# add_service_files(
#   FILES
#   Service1.srv
#   Service2.srv
# )

## Generate added messages and services with any dependencies listed here
# generate_messages(
#   DEPENDENCIES
#   std_msgs  # Or other packages containing msgs
# )

###################################
## catkin specific configuration ##
###################################
## The catkin_package macro generates cmake config files for your package
## Declare things to be passed to dependent projects
## LIBRARIES: libraries you create in this project that dependent projects also need
## CATKIN_DEPENDS: catkin_packages dependent projects also need
## DEPENDS: system dependencies of this project that dependent projects also need
catkin_package(
#  INCLUDE_DIRS include
#  LIBRARIES rosgo
#  CATKIN_DEPENDS other_catkin_pkg
#  DEPENDS system_lib
)


###########
## Build ##
###########

## Specify additional locations of header files
## Your package locations should be listed before other locations
# include_directories(include ${catkin_INCLUDE_DIRS} ${Boost_INCLUDE_DIRS})

## Declare a cpp library
# add_library(rosgo
#   src/${PROJECT_NAME}/rosgo.cpp
# )

## Declare a cpp executable
# add_executable(rosgo_node src/rosgo_node.cpp)

## Add cmake target dependencies of the executable/library
## as an example, message headers may need to be generated before nodes
# add_dependencies(rosgo_node rosgo_generate_messages_cpp)

## Specify libraries to link a library or executable target against
# target_link_libraries(rosgo_node
#   ${catkin_LIBRARIES}
# )

#set(ROSGO_ROOT "${CATKIN_DEVEL_PREFIX}/go")
#file(MAKE_DIRECTORY ${ROSGO_ROOT})
#
#set(ROSGO_PATH "${PROJECT_SOURCE_DIR}:${ROSGO_ROOT}")
#message(STATUS "ROSGO_PATH=${ROSGO_PATH}")
#
#set(ROSGO_SRC ${ROSGO_ROOT}/src)
#file(MAKE_DIRECTORY ${ROSGO_SRC})
#
#set(ROSGO_BIN ${ROSGO_ROOT}/bin)
#file(MAKE_DIRECTORY ${ROSGO_BIN})
#
#execute_process(
#    COMMAND go env GOARCH
#    OUTPUT_VARIABLE ROSGO_ARCH
#    OUTPUT_STRIP_TRAILING_WHITESPACE)
#execute_process(
#    COMMAND go env GOOS
#    OUTPUT_VARIABLE ROSGO_OS
#    OUTPUT_STRIP_TRAILING_WHITESPACE)
#set(ROSGO_PKG ${ROSGO_ROOT}/pkg/${ROSGO_OS}_${ROSGO_ARCH})
#file(MAKE_DIRECTORY ${ROSGO_PKG})
#
#
#set(xmlrpc ${PROJECT_SOURCE_DIR}/pkg/${ROSGO_OS}_${ROSGO_ARCH}/xmlrpc.a)
#add_custom_command(
#    OUTPUT ${xmlrpc}
#    COMMAND env GOPATH="${PROEJCT_SOURCE_DIR}" go install xmlrpc)
#add_custom_target(${PROJECT_NAME}_xmlrpc ALL DEPENDS ${xmlrpc})
#
#set(rosgo ${PROJECT_SOURCE_DIR}/pkg/${ROSGO_OS}_${ROSGO_ARCH}/ros.a)
#add_custom_command(
#    OUTPUT ${rosgo} 
#    COMMAND env GOPATH="${PROJECT_SOURCE_DIR}" go install ros
#    DEPENDS ${PROJECT_NAME}_xmlrpc)
#add_custom_target(${PROJECT_NAME}_ros ALL DEPENDS ${rosgo})
#
#if(ENV{GOPATH})
#    set(ENV{GOPATH} "${PROJECT_SOURCE_DIR}:${ROSGO_ROOT}:ENV{GOPATH}")
#else()
#    set(ENV{GOPATH} "${PROJECT_SOURCE_DIR}:${ROSGO_ROOT}")
#endif()


rosgo_force_generate_messages(std_msgs rosgraph_msgs)
rosgo_add_library(xmlrpc TARGET ${PROJECT_NAME}_xmlrpc)
rosgo_add_library(
        ros
        TARGET ${PROJECT_NAME}
        DEPENDS ${PROJECT_NAME}_xmlrpc std_msgs_gengo rosgraph_msgs_gengo
    )

rosgo_add_test(xmlrpc DEPENDS ${PROJECT_NAME}_xmlrpc)
rosgo_add_test(ros DEPENDS ${PROJECT_NAME})

#############
## Install ##
#############

# all install targets should use catkin DESTINATION variables
# See http://ros.org/doc/groovy/api/catkin/html/adv_user_guide/variables.html

## Mark executable scripts (Python etc.) for installation
## in contrast to setup.py, you can choose the destination
# install(PROGRAMS
#   scripts/my_python_script
#   DESTINATION ${CATKIN_PACKAGE_BIN_DESTINATION}
# )

## Mark executables and/or libraries for installation
# install(TARGETS rosgo rosgo_node
#   ARCHIVE DESTINATION ${CATKIN_PACKAGE_LIB_DESTINATION}
#   LIBRARY DESTINATION ${CATKIN_PACKAGE_LIB_DESTINATION}
#   RUNTIME DESTINATION ${CATKIN_PACKAGE_BIN_DESTINATION}
# )

## Mark cpp header files for installation
# install(DIRECTORY include/${PROJECT_NAME}/
#   DESTINATION ${CATKIN_PACKAGE_INCLUDE_DESTINATION}
#   FILES_MATCHING PATTERN "*.h"
#   PATTERN ".svn" EXCLUDE
# )

## Mark other files for installation (e.g. launch and bag files, etc.)
# install(FILES
#   # myfile1
#   # myfile2
#   DESTINATION ${CATKIN_PACKAGE_SHARE_DESTINATION}
# )


#############
## Testing ##
#############

## Add gtest based cpp test target and link libraries
# catkin_add_gtest(${PROJECT_NAME}-test test/test_rosgo.cpp)
# if(TARGET ${PROJECT_NAME}-test)
#   target_link_libraries(${PROJECT_NAME}-test ${PROJECT_NAME})
# endif()

## Add folders to be run by python nosetests
# catkin_add_nosetests(test)
