.PHONY: checks
checks:
	@# 1. Configure (if not already done)
	cmake -S . -B build -DCMAKE_BUILD_TYPE=Debug -DCMAKE_EXPORT_COMPILE_COMMANDS=ON
	@# 2. Build
	cmake --build build
	@# 3. Run tests
	ctest --test-dir build --output-on-failure
