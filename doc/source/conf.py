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

sys.path.insert(0, os.path.abspath("../../"))

from ansys_sphinx_theme import ansys_favicon

# Static version - no VERSION file dependency
version_file = "1.0.0"

# Project metadata
project = "AALI SharedTypes"
author = "ANSYS, Inc."
copyright = f"{datetime.now().year} ANSYS, Inc. All rights reserved"

release = version = version_file
switcher_version = version_file
cname = os.getenv("DOCUMENTATION_CNAME", "noname.com")

extensions = [
    "sphinx.ext.autodoc",
    "sphinx.ext.napoleon",
    "sphinx.ext.viewcode",
    "sphinx.ext.graphviz",
    "sphinx_design",
    "sphinx_external_toc",
]

html_theme = "ansys_sphinx_theme"
html_favicon = ansys_favicon
html_short_title = html_title = project

external_toc_path = "_toc.yml"
external_toc_exclude_missing = False

html_context = {
    "github_user": "ansys",
    "github_repo": "aali-sharedtypes",
    "github_version": "main",
    "doc_path": "doc/source",
}
html_theme_options = {
    "logo": "ansys",
    "github_url": "https://github.com/ansys/aali-sharedtypes",
    "additional_breadcrumbs": [
        ("AALI", "https://aali.docs.pyansys.com/"),
    ],
    "switcher": {
        "json_url": "_static/versions.json",
        "version_match": switcher_version,
    },
    "navbar_end": ["navbar-icon-links", "version-switcher", "theme-switcher"],
    "check_switcher": True,
    "show_prev_next": True,
    "show_breadcrumbs": True,
    "use_edit_page_button": True,
    "navigation_depth": 4,
    "collapse_navigation": False,
}

html_css_files = [
    "https://fonts.googleapis.com/icon?family=Material+Icons",
    "custom.css"
]

html_static_path = ["_static"]
