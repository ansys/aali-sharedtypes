

import os
import sys
from datetime import datetime

# Project metadata
project = "Aali Workflow Config"
author = "Aali Team"
copyright = f"{datetime.now().year} ANSYS, Inc"

sys.path.insert(0, os.path.abspath("../../"))

extensions = [

    "sphinx.ext.autodoc",
    "sphinx.ext.napoleon",
    "sphinx.ext.viewcode",
    "sphinx_design",
    "sphinx_external_toc",
]

html_theme = "ansys_sphinx_theme"

external_toc_path = "_toc.yml"
external_toc_exclude_missing = False


html_css_files = ["https://fonts.googleapis.com/icon?family=Material+Icons"]
