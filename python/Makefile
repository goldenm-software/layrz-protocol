.PHONY: checks
checks:
	uv sync --group dev
	uv run ruff check
	uv run ty check
	uv run coverage run -m pytest tests/ -v
	uv run coverage report -m

.PHONY: coverage
coverage:
	uv sync --group dev
	mkdir -p coverage
	uv run coverage run -m pytest tests/ -v
	uv run coverage lcov -o coverage/lcov.info
	uv run coverage report -m
