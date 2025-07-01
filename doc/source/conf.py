# Copyright (C) 2025 ANSYS, Inc. and/or its affiliates.
# SPDX-License-Identifier: MIT
#
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.


import os
import sys
from datetime import datetime


sys.path.insert(0, os.path.abspath("../../pkg"))

# Project metadata
project = "Aali Workflow Config"
author = "Aali Team"
copyright = f"{datetime.now().year} ANSYS, Inc"

# Ensure source code is importable
sys.path.insert(0, os.path.abspath("../../"))

# Sphinx extensions
extensions = [
    "sphinx.ext.autodoc",        # Enables auto API generation
    "sphinx.ext.napoleon",       # Supports Google/NumPy-style docstrings
    "sphinx.ext.autosummary",    # Generates summary tables
    "sphinx.ext.viewcode",       # Adds links to source code
    "sphinx_design",             # Enables layout tools like grids/cards
    "sphinx_external_toc",       # External _toc.yml support
]

# Automatically generate stub pages from autosummary directives
autosummary_generate = True

# HTML theme configuration
html_theme = "ansys_sphinx_theme"

# External table of contents (_toc.yml)
external_toc_path = "_toc.yml"
external_toc_exclude_missing = False

# Add custom CSS files (Google Material icons)
html_css_files = [
    "https://fonts.googleapis.com/icon?family=Material+Icons"
]
