# ── per-module Makefiles ────────────────────────────────────────────────────
# Each module (agni-router, agni-agent, agni-nova) has its own Makefile with
# platform-specific targets. Use the targets below as convenience shortcuts
# from the repo root, or cd into the module directory and run make directly.
#
# Output binaries are placed in release/ at the project root.
# ─────────────────────────────────────────────────────────────────────────────

RELEASE_DIR := release

# ── agent ────────────────────────────────────────────────────────────────────
agent-linux:
	$(MAKE) -C agni-agent build-linux

agent-darwin:
	$(MAKE) -C agni-agent build-darwin

agent-windows:
	$(MAKE) -C agni-agent build-windows

agent-all:
	$(MAKE) -C agni-agent build-all


## Show this help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Agent targets (agni-agent/Makefile):"
	@echo "  agent-linux      Build agni-agent for Linux"
	@echo "  agent-darwin     Build agni-agent for macOS"
	@echo "  agent-windows    Build agni-agent for Windows"
	@echo "  agent-all        Build agni-agent for all platforms"
	@echo ""
	@echo "All binaries are written to: $(RELEASE_DIR)/{linux,darwin,windows}/"

.PHONY: agent-linux  agent-darwin  agent-windows  agent-all  help