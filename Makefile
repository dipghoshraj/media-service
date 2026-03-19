# ── per-module Makefiles ────────────────────────────────────────────────────
# Each module (agni-router, agni-agent, agni-nova) has its own Makefile with
# platform-specific targets. Use the targets below as convenience shortcuts
# from the repo root, or cd into the module directory and run make directly.
#
# Output binaries are placed in release/ at the project root.
# ─────────────────────────────────────────────────────────────────────────────

RELEASE_DIR := release

# ── router ───────────────────────────────────────────────────────────────────
router-linux:
	$(MAKE) -C agni-router build-linux

router-darwin:
	$(MAKE) -C agni-router build-darwin

router-windows:
	$(MAKE) -C agni-router build-windows

router-all:
	$(MAKE) -C agni-router build-all

router-certs:
	$(MAKE) -C agni-router gen-certs

# ── agent ────────────────────────────────────────────────────────────────────
agent-linux:
	$(MAKE) -C agni-agent build-linux

agent-darwin:
	$(MAKE) -C agni-agent build-darwin

agent-windows:
	$(MAKE) -C agni-agent build-windows

agent-all:
	$(MAKE) -C agni-agent build-all

# ── nova ─────────────────────────────────────────────────────────────────────
nova-linux:
	$(MAKE) -C agni-nova build-linux

nova-darwin:
	$(MAKE) -C agni-nova build-darwin

nova-windows:
	$(MAKE) -C agni-nova build-windows

nova-all:
	$(MAKE) -C agni-nova build-all

# ── combined ─────────────────────────────────────────────────────────────────
build-all-linux: router-linux agent-linux nova-linux

build-all-darwin: router-darwin agent-darwin nova-darwin

build-all-windows: router-windows agent-windows nova-windows

build-all: router-all agent-all nova-all

## Show this help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Router targets (agni-router/Makefile):"
	@echo "  router-linux     Build agni-router for Linux"
	@echo "  router-darwin    Build agni-router for macOS"
	@echo "  router-windows   Build agni-router for Windows"
	@echo "  router-all       Build agni-router for all platforms"
	@echo "  router-certs     Generate router TLS certificates"
	@echo ""
	@echo "Agent targets (agni-agent/Makefile):"
	@echo "  agent-linux      Build agni-agent for Linux"
	@echo "  agent-darwin     Build agni-agent for macOS"
	@echo "  agent-windows    Build agni-agent for Windows"
	@echo "  agent-all        Build agni-agent for all platforms"
	@echo ""
	@echo "Nova targets (agni-nova/Makefile):"
	@echo "  nova-linux       Build agni-nova for Linux"
	@echo "  nova-darwin      Build agni-nova for macOS"
	@echo "  nova-windows     Build agni-nova for Windows"
	@echo "  nova-all         Build agni-nova for all platforms"
	@echo ""
	@echo "Combined targets:"
	@echo "  build-all-linux   Build all components for Linux"
	@echo "  build-all-darwin  Build all components for macOS"
	@echo "  build-all-windows Build all components for Windows"
	@echo "  build-all         Build all components for all platforms"
	@echo ""
	@echo "All binaries are written to: $(RELEASE_DIR)/"

.PHONY: router-linux router-darwin router-windows router-all router-certs \
        agent-linux  agent-darwin  agent-windows  agent-all  \
        nova-linux   nova-darwin   nova-windows   nova-all   \
        build-all-linux build-all-darwin build-all-windows build-all help