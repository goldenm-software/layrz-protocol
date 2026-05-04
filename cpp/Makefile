.PHONY: checks
checks:
	@# 1. Configure (if not already done)
	cmake -S . -B build -DCMAKE_BUILD_TYPE=Debug -DCMAKE_EXPORT_COMPILE_COMMANDS=ON
	@# 2. Build
	cmake --build build
	@# 3. Run tests
	ctest --test-dir build --output-on-failure

.PHONY: coverage
coverage:
	@# 1. Configure with coverage flags
	cmake -S . -B build-cov -DCMAKE_BUILD_TYPE=Debug \
	      -DCMAKE_CXX_FLAGS="--coverage -O0" \
	      -DCMAKE_EXE_LINKER_FLAGS="--coverage"
	@# 2. Build
	cmake --build build-cov
	@# 3. Run tests
	ctest --test-dir build-cov --output-on-failure
	@# 4. Capture coverage
	@mkdir -p coverage
	lcov --capture --directory build-cov --output-file coverage/lcov.info \
	     --ignore-errors mismatch,gcov,unused
	lcov --remove coverage/lcov.info '/usr/*' '*/build-cov/_deps/*' \
	     -o coverage/lcov.info
