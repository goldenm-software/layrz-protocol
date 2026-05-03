cmake_minimum_required(VERSION 3.16)
project(layrz_protocol VERSION 3.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

option(LAYRZ_PROTOCOL_BUILD_NET   "Build HTTP and TCP transport layer (requires pthreads)" ON)
option(LAYRZ_PROTOCOL_BUILD_TESTS "Build tests (requires GoogleTest via FetchContent)"    ON)

# ─── Core library (Tier 1) ───────────────────────────────────────────────────
file(GLOB_RECURSE CORE_SOURCES
    src/crc.cpp
    src/float_repr.cpp
    src/extras.cpp
    src/parser.cpp
    src/ble/*.cpp
    src/packets/*.cpp
)

add_library(layrz_protocol_core ${CORE_SOURCES})
add_library(layrz::protocol::core ALIAS layrz_protocol_core)

target_include_directories(layrz_protocol_core
    PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/src
)

target_compile_features(layrz_protocol_core PUBLIC cxx_std_17)

# ─── Net library (Tier 2) ─────────────────────────────────────────────────────
if(LAYRZ_PROTOCOL_BUILD_NET)
    find_package(Threads REQUIRED)

    # Prefer system cpp-httplib; fall back to FetchContent if not found
    find_path(HTTPLIB_INCLUDE_DIR httplib.h)
    if(NOT HTTPLIB_INCLUDE_DIR)
        include(FetchContent)
        FetchContent_Declare(
            cpp-httplib
            GIT_REPOSITORY https://github.com/yhirose/cpp-httplib.git
            GIT_TAG        v0.18.3
        )
        FetchContent_MakeAvailable(cpp-httplib)
        set(HTTPLIB_INCLUDE_DIR "${cpp-httplib_SOURCE_DIR}")
    endif()

    add_library(layrz_protocol_net
        src/transport/http_comm.cpp
        src/transport/tcp_comm.cpp
    )
    add_library(layrz::protocol::net ALIAS layrz_protocol_net)

    target_include_directories(layrz_protocol_net
        PUBLIC
            $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
            $<INSTALL_INTERFACE:include>
        PRIVATE
            ${HTTPLIB_INCLUDE_DIR}
    )

    target_link_libraries(layrz_protocol_net
        PUBLIC  layrz_protocol_core
        PRIVATE Threads::Threads
    )

    target_compile_definitions(layrz_protocol_net
        PUBLIC  LAYRZ_PROTOCOL_CLIENTS
        PRIVATE CPPHTTPLIB_OPENSSL_SUPPORT=0
    )
endif()

# ─── Tests ────────────────────────────────────────────────────────────────────
if(LAYRZ_PROTOCOL_BUILD_TESTS)
    include(FetchContent)
    FetchContent_Declare(
        googletest
        GIT_REPOSITORY https://github.com/google/googletest.git
        GIT_TAG        v1.15.2
    )
    set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
    FetchContent_MakeAvailable(googletest)

    enable_testing()
    add_subdirectory(tests)
endif()

# ─── Examples ─────────────────────────────────────────────────────────────────
add_executable(example_encode_decode examples/encode_decode.cpp)
target_link_libraries(example_encode_decode PRIVATE layrz_protocol_core)

if(LAYRZ_PROTOCOL_BUILD_NET)
    add_executable(example_tcp_client examples/tcp_client.cpp)
    target_link_libraries(example_tcp_client PRIVATE layrz_protocol_net)
endif()

# ─── Install ──────────────────────────────────────────────────────────────────
include(GNUInstallDirs)
include(CMakePackageConfigHelpers)

install(TARGETS layrz_protocol_core
    EXPORT  LayrzProtocolTargets
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
    ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
)

if(LAYRZ_PROTOCOL_BUILD_NET)
    install(TARGETS layrz_protocol_net
        EXPORT  LayrzProtocolTargets
        LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
        ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
    )
endif()

install(DIRECTORY include/ DESTINATION ${CMAKE_INSTALL_INCLUDEDIR})

install(EXPORT LayrzProtocolTargets
    FILE      LayrzProtocolTargets.cmake
    NAMESPACE layrz::protocol::
    DESTINATION ${CMAKE_INSTALL_LIBDIR}/cmake/LayrzProtocol
)

configure_package_config_file(
    cmake/LayrzProtocolConfig.cmake.in
    ${CMAKE_CURRENT_BINARY_DIR}/LayrzProtocolConfig.cmake
    INSTALL_DESTINATION ${CMAKE_INSTALL_LIBDIR}/cmake/LayrzProtocol
)

install(FILES ${CMAKE_CURRENT_BINARY_DIR}/LayrzProtocolConfig.cmake
    DESTINATION ${CMAKE_INSTALL_LIBDIR}/cmake/LayrzProtocol
)
