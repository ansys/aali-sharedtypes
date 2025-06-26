# Minimal Makefile for building Sphinx documentation

SPHINXBUILD   = sphinx-build
SOURCEDIR     = source
BUILDDIR      = _build

.PHONY: help clean html

help:
	@$(SPHINXBUILD) -M help "$(SOURCEDIR)" "$(BUILDDIR)"

clean:
	rm -rf $(BUILDDIR)

html:
	@$(SPHINXBUILD) -M html "$(SOURCEDIR)" "$(BUILDDIR)"