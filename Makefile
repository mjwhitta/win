-include gomk/main.mk
-include local/Makefile

ifneq ($(unameS),windows)
spellcheck:
	@codespell -f -S ".git,generated.go"
endif
