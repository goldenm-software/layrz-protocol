.PHONY: checks
checks:
	uv sync --group dev
	uv run ruff check
	uv run mypy .
	uv run coverage run -m pytest tests/ -v
	uv run coverage report -m
