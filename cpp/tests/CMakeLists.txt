include(GoogleTest)

# Helper to add a test target linked to core + gtest
function(add_layrz_test name)
    add_executable(${name} ${name}.cpp)
    target_link_libraries(${name}
        PRIVATE
            layrz_protocol_core
            GTest::gtest_main
    )
    target_include_directories(${name}
        PRIVATE
            ${CMAKE_CURRENT_SOURCE_DIR}
    )
    gtest_discover_tests(${name})
endfunction()

add_layrz_test(crc_test)
add_layrz_test(float_repr_test)
add_layrz_test(extras_test)
add_layrz_test(roundtrip_test)
add_layrz_test(ble_test)
