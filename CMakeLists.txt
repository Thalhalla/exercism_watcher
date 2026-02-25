# Minimum CMake version required
cmake_minimum_required(VERSION 3.14)

project(exercism_watcher)

# Specify the script file to be installed
set(EXERSCISM_WATCHER_NAME "exercism_watcher")
# Path to Go source file
set(EXERSCISM_WATCHER_SRC "${CMAKE_CURRENT_SOURCE_DIR}/${EXERSCISM_WATCHER_NAME}.go")
# Output binary path
set(EXERSCISM_WATCHER_BIN "${CMAKE_CURRENT_BINARY_DIR}/${EXERSCISM_WATCHER_NAME}")
# Custom command to build the Go binary
add_custom_command(
    OUTPUT ${EXERSCISM_WATCHER_BIN}
    COMMAND go build -o ${EXERSCISM_WATCHER_BIN} ${EXERSCISM_WATCHER_SRC}
    DEPENDS ${EXERSCISM_WATCHER_SRC}
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    COMMENT "Building Go executable ${EXERSCISM_WATCHER_NAME}"
)
# Custom target that runs the build
add_custom_target(build_exercism_watcher ALL DEPENDS ${EXERSCISM_WATCHER_BIN})
# Install the built binary instead of the source script
set(EXERSCISM_WATCHER_PATH "${EXERSCISM_WATCHER_BIN}")

# Install the script to a destination directory
# The 'PROGRAMS' form ensures executable permissions are set.
# 'DESTINATION bin' places it in the 'bin' subdirectory of the install prefix.
install(PROGRAMS ${EXERSCISM_WATCHER_PATH} 
    DESTINATION bin
    PERMISSIONS OWNER_READ OWNER_WRITE OWNER_EXECUTE GROUP_READ GROUP_EXECUTE WORLD_READ WORLD_EXECUTE
)

# Install the man page with proper section handling
install(
    FILES ${CMAKE_CURRENT_SOURCE_DIR}/man/exercism_watcher.1
    DESTINATION share/man/man1
    COMPONENT Runtime
)
# --------------------------------------------------------------
# Optional: make the install prefix visible to the user (helpful for packaging)
# --------------------------------------------------------------
message(STATUS "Install prefix: ${CMAKE_INSTALL_PREFIX}")
message(STATUS "Binary will be installed to: ${CMAKE_INSTALL_PREFIX}/bin")
message(STATUS "Man page will be installed to: ${CMAKE_INSTALL_PREFIX}/share/man/man1")
