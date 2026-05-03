CHMOD_CMD = chmod +x .githooks/pre-commit
ifeq ($(OS),Windows_NT)
    CHMOD_CMD = echo "Skipping chmod on Windows"
endif

.PHONY: install-hooks
install-hooks:
	@echo "Installing git hooks from .githooks directory..."
	@$(CHMOD_CMD)
	@git config core.hooksPath .githooks

.PHONY: checks
checks:
	@echo "Running checks for all languages..."
	@$(MAKE) -C python checks
	@$(MAKE) -C dart checks
	@$(MAKE) -C go checks
