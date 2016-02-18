# vim: ft=cmake :

function(_rosgo_setup_global_variable)
    set(libdir "${CATKIN_DEVEL_PREFIX}/lib")
    set(root "${libdir}/go")
    file(MAKE_DIRECTORY ${root})
    execute_process(COMMAND go env GOARCH OUTPUT_VARIABLE goarch OUTPUT_STRIP_TRAILING_WHITESPACE)
    execute_process(COMMAND go env GOOS OUTPUT_VARIABLE goos OUTPUT_STRIP_TRAILING_WHITESPACE)
    set_property(GLOBAL PROPERTY _ROSGO_ROOT "${root}")
    set_property(GLOBAL PROPERTY _ROSGO_BIN "${libdir}")
    set_property(GLOBAL PROPERTY _ROSGO_SRC "${root}/src")
    set_property(GLOBAL PROPERTY _ROSGO_PKG "${root}/pkg/${goos}_${goarch}")
    get_property(gopath GLOBAL PROPERTY _ROSGO_PATH)
    if("${gopath}" STREQUAL "")
        set_property(GLOBAL PROPERTY _ROSGO_PATH "${root}")
    endif()
    set_property(GLOBAL APPEND PROPERTY _ROSGO_PATH "${PROJECT_SOURCE_DIR}")

    if("${catkin_GO_LIBRARIES}" STREQUAL "")
        set(catkin_GO_LIBRARIES "" PARENT_SCOPE)
    endif()
endfunction()


# This will be evaluated per each project that depend `rosgo_build_tools`.
_rosgo_setup_global_variable()


function(_rosgo_make_gopath gopath_result)
  if("$ENV{GOPATH}" STREQUAL "")
    set(${gopath_result} ${CATKIN_DEVEL_PREFIX}/lib/go PARENT_SCOPE)
  else()
    set(${gopath_result} ${CATKIN_DEVEL_PREFIX}/lib/go:$ENV{GOPATH} PARENT_SCOPE)
  endif()
endfunction()


# Clear old symlinks and create new ones that point original sources.
function(_rosgo_mirror_go_files package var)
    get_filename_component(orig_dir "${PROJECT_SOURCE_DIR}/src/${package}" ABSOLUTE)
    get_filename_component(link_dir "${CATKIN_DEVEL_PREFIX}/lib/go/src/${package}" ABSOLUTE)

    file(MAKE_DIRECTORY "${link_dir}")

    file(GLOB orig_paths "${orig_dir}/*.go")
    set(filenames "")
    foreach(p ${orig_paths})
        get_filename_component(f ${p} NAME)
        list(APPEND filenames "${f}")
    endforeach()

    file(GLOB last_items "${link_dirs}/*.go")
    foreach(item ${last_items})
        if(IS_SYMLINK "${item}")
            file(REMOVE "${item}")
        endif()
    endforeach()

    set(links "")
    foreach(filename ${filenames})
        set(orig "${orig_dir}/${filename}")
        set(link "${link_dir}/${filename}")
        add_custom_command(
            OUTPUT "${link}"
            COMMAND ${CMAKE_COMMAND} -E create_symlink "${orig}" "${link}"
            )
        list(APPEND links "${link}")
    endforeach()
    set(${var} ${links} PARENT_SCOPE)
endfunction()


# Add executable target
function(catkin_add_go_executable)
    set(options)
    set(one_value_args TARGET)
    set(multi_value_args DEPENDS)
    cmake_parse_arguments(catkin_add_go_executable "${options}" "${one_value_args}"
                          "${multi_value_args}" "${ARGN}")
    list(GET catkin_add_go_executable_UNPARSED_ARGUMENTS 0 package)
    if("${catkin_add_go_executable_TARGET}" STREQUAL "")
        string(REPLACE "/" "_" target "${PROJECT_NAME}_${package}")
        if(NOT ${target} STREQUAL ${PROJECT_NAME}_NOTFOUND)
            set(catkin_add_go_executable_TARGET ${target})
        endif()
    endif()

    _rosgo_mirror_go_files(${package} src_links)

    string(REPLACE "/" ";" exe_path "${package}")
    list(GET exe_path -1 exe_name)
    set(exe "${CATKIN_DEVEL_PREFIX}/lib/${PROJECT_NAME}/${exe_name}")

    _rosgo_make_gopath(gopath)

    add_custom_target(
            ${catkin_add_go_executable_TARGET} ALL
            COMMAND env GOPATH=${gopath} go build -o ${exe} ${package}
            DEPENDS ${catkin_add_go_executable_DEPENDS} ${src_links})
endfunction()


# Add library target
function(catkin_add_go_library)
    set(options)
    set(one_value_args TARGET)
    set(multi_value_args DEPENDS)
    cmake_parse_arguments(catkin_add_go_library "${options}" "${one_value_args}"
                          "${multi_value_args}" "${ARGN}")
    list(GET catkin_add_go_library_UNPARSED_ARGUMENTS 0 package)
    if("${catkin_add_go_library_TARGET}" STREQUAL "")
        string(REPLACE "/" "_" target "${PROJECT_NAME}_${package}")
        if(NOT ${target} STREQUAL ${PROJECT_NAME}_NOTFOUND)
            set(catkin_add_go_library_TARGET ${target})
        endif()
    endif()

    _rosgo_mirror_go_files(${package} src_links)
    get_property(gopkg GLOBAL PROPERTY _ROSGO_PKG)

    _rosgo_make_gopath(gopath)

    add_custom_target(
            ${catkin_add_go_library_TARGET} ALL
            COMMAND env GOPATH=${gopath} go build -o ${gopkg}/${package}.a ${package}
            DEPENDS ${catkin_add_go_library_DEPENDS} ${src_links})
    list(INSERT catkin_GO_LIBRARIES 0 ${catkin_add_go_library_TARGET})
endfunction()


# Add test target
function(catkin_add_go_test)
    set(options)
    set(one_value_args)
    set(multi_value_args DEPENDS)
    cmake_parse_arguments(catkin_add_go_test "${options}" "${one_value_args}"
                          "${multi_value_args}" "${ARGN}")
    list(GET catkin_add_go_test_UNPARSED_ARGUMENTS 0 package)
    string(REPLACE "/" "_" target "${package}")

    _rosgo_mirror_go_files(${package} src_links)

    _rosgo_make_gopath(gopath)

    set(_depends "${catkin_add_go_test_DEPENDS};${src_links}")

    add_custom_target(
        run_tests_${PROJECT_NAME}_go_test_${target}
        #COMMAND env GOPATH=${gopath} go test ${package}
        COMMAND ${CATKIN_ENV} env GOPATH=${gopath} rosrun rosgo_build_tools rosgo-test-wrapper.sh ${package}
        DEPENDS ${_depends})

    # Register this test to workspace-wide run_tests target
    if(NOT TARGET run_tests_${PROJECT_NAME}_go_test)
        add_custom_target(run_tests_${PROJECT_NAME}_go_test)
        add_dependencies(run_tests run_tests_${PROJECT_NAME}_go_test)
    endif()
    add_dependencies(run_tests_${PROJECT_NAME}_go_test
                     run_tests_${PROJECT_NAME}_go_test_${target})
endfunction()


