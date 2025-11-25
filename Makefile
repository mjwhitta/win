-include gomk/main.mk
-include local/Makefile

CC := x86_64-w64-mingw32-gcc
CGO_ENABLED := 1
GOARCH := amd64
GOOS := windows

ifneq ($(unameS),windows)
spellcheck:
	@codespell -f -L hilighter -S "*.pem,.git,generated.go,go.*,gomk"
endif
