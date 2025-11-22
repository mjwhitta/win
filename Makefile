-include gomk/main.mk
-include local/Makefile

ifneq ($(unameS),windows)
spellcheck:
	@codespell -f -L hilighter -S "*.pem,.git,generated.go,go.*,gomk"
endif
