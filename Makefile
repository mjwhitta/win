-include gomk/main.mk
-include local/Makefile

CC := x86_64-w64-mingw32-gcc
GOARCH := amd64
GOOS := windows
OUT := $(BUILD)/$(GOOS)/$(GOARCH)

clean: clean-default
ifeq ($(unameS),windows)
ifneq ($(wildcard resource_windows*.syso),)
	@remove-item -force ./cmd/*/resource_windows*.syso
endif
else
	@rm -f ./cmd/*/resource_windows*.syso
endif

ifneq ($(unameS),windows)
spellcheck:
	@codespell -f -L hilighter -S "*.pem,.git,generated.go,go.*,gomk"
endif
