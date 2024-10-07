cmake_minimum_required(VERSION 3.20)
project(SpectreLink DESCRIPTION "SpectreLink" VERSION 0.1.0 LANGUAGES C)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED OFF)
set(CMAKE_C_EXTENSIONS OFF)

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wall -Werror -Wextra -Wno-unused-parameter -Wno-unused-function -Wno-unused-variable -Wno-missing-field-initializers -Wno-missing-braces -Wno-unused-result")
add_executable(
        ${PROJECT_NAME}
        src/main.c
        src/utils/socket.c
        src/utils/urlparser.c
        src/utils/spectre_time.c
        src/utils/socket.c
        src/utils/debug.c
        src/utils/common.c
        src/protocol/bot_info.c
        src/protocol/connection.c
        src/protocol/packet.c
        src/protocol/protocol.c
)

include_directories(${PROJECT_NAME} PRIVATE
        src/protocol/includes
        src/bot/includes
        src/utils/includes
)

add_definitions(-DBOT_VERSION="${CMAKE_PROJECT_VERSION}")

if (WIN32)
    target_link_libraries(${PROJECT_NAME} ws2_32)
endif()

if ("${CMAKE_BUILD_TYPE}" STREQUAL "Debug")
    target_compile_definitions(${PROJECT_NAME} PRIVATE -DDEBUG)
endif()