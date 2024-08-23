BINDIR      := $(CURDIR)/bin
BINNAME     ?= helm-dev
HELM_PLUGIN_DIR ?= $(shell helm env HELM_PLUGINS)

.PHONY: build
build:
	go build -o $(BINDIR)/$(BINNAME) ./cmd/helm-dev

# Create a symlink to the plugins director so that the plugin can be used
.PHONY: link
link:
	mkdir -p $(HELM_PLUGIN_DIR)
	ln -s $(CURDIR) $(HELM_PLUGIN_DIR)/$(BINNAME)