# Copyright (C) 2025 - 2026 ANSYS, Inc. and/or its affiliates.
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
import shutil
import re
import sys
import pathlib

# Ensure the correct number of arguments are provided
if len(sys.argv) != 3:
    print("Usage: python script.py <repo_owner> <repo_name>")
    sys.exit(1)

# Parse CLI arguments
repo_owner = sys.argv[1]
repo_name = sys.argv[2]

# Set environment variables for file paths
DOC_BUILD_HTML = "documentation-html"
SOURCE_FILE = os.path.join(DOC_BUILD_HTML, "api_reference", "test", "index.html")
SOURCE_DIRECTORY = f"dist/pkg/github.com/{repo_owner}/{repo_name}/pkg/"
REPLACEMENT_DIRECTORY = os.path.join(DOC_BUILD_HTML, "api_reference", "pkg")
ACTUAL_DIR = os.path.join(DOC_BUILD_HTML, "api_reference/")
STATIC_DIR = os.path.join(DOC_BUILD_HTML, "_static")


def get_relative_path_to_static_files(file_path):
    # Get the relative path to the static files
    relative_path = os.path.relpath(STATIC_DIR, start=file_path)
    return relative_path

# Check if REPLACEMENT_DIRECTORY exists, if not, create it
if not os.path.exists(REPLACEMENT_DIRECTORY):
    os.makedirs(REPLACEMENT_DIRECTORY)

# Remove existing content in the replacement directory
for filename in os.listdir(REPLACEMENT_DIRECTORY):
    file_path = os.path.join(REPLACEMENT_DIRECTORY, filename)
    try:
        if os.path.isfile(file_path) or os.path.islink(file_path):
            os.unlink(file_path)
        elif os.path.isdir(file_path):
            shutil.rmtree(file_path)
    except Exception as e:
        print(f'Failed to delete {file_path}. Reason: {e}')

# Move the source_directory content to the replacement_directory
for filename in os.listdir(SOURCE_DIRECTORY):
    shutil.move(os.path.join(SOURCE_DIRECTORY, filename), REPLACEMENT_DIRECTORY)

# Remove the index.html file in the replacement_directory
index_file_path = os.path.join(REPLACEMENT_DIRECTORY, "index.html")
if os.path.exists(index_file_path):
    os.remove(index_file_path)

def get_subfolder_path(file_path, source_directory):
    """Get the relative path of the file to the source directory.

    Parameters
    ----------
    file_path : str
        The path of the file.
    source_directory : str
        The path of the source directory.

    Returns
    -------
    str
        The relative path of the file to the source directory.
    """
    relative_path = os.path.relpath(file_path, start=source_directory)
    return relative_path


# Process each HTML file in the replacement directory
for root, dirs, files in os.walk(REPLACEMENT_DIRECTORY):
    for file in files:
        if file.endswith(".html"):
            replacement_file = os.path.join(root, file)
            with open(replacement_file, 'r') as f:
                content = f.read()

            # Extract body content and perform replacements
            match = re.search(r'<body>([\s\S]*?)<\/body>', content)
            if match:
                replacementBodyContent = match.group(1)
                search_pattern = (
                    r'<div class="top-heading" id="heading-wide"><a href="/pkg/github.com/'
                    + re.escape(repo_owner)
                    + '/'
                    + re.escape(repo_name)
                    + r'/">'
                    + re.escape("GoPages | Auto-generated docs")
                    + r'</a></div>'
                    + r'[\s\S]*?<a href="#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>'
                )
                replacementBodyContent = re.sub(search_pattern, '', replacementBodyContent)

                # Read source file and replace content
                with open(SOURCE_FILE, 'r') as f:
                    source_content = f.read()

                escaped_replacement = replacementBodyContent.replace("\\", "\\\\")

                new_content = re.sub(
                    r'(<article class="bd-article">)[\s\S]*?(<\/article>)',
                    rf'\1{escaped_replacement}\2',
                    source_content
                )

                # Write the modified content back to the replacement file
                with open(replacement_file, 'w') as f:
                    f.write(new_content)
            else:
                print(f'No <body> content found in {replacement_file}')

# Move the modified files back to the actual directory
for filename in os.listdir(REPLACEMENT_DIRECTORY):
    shutil.move(os.path.join(REPLACEMENT_DIRECTORY, filename), ACTUAL_DIR)

# Find the subfolders in the actual directory and fix the static paths
def find_subfolders():
    """Find the subfolders in the actual directory and return the list of subfolders.

    Returns
    -------
    list
        The list of subfolders in the actual directory.
    """
    for root, dirs, files in os.walk(ACTUAL_DIR):
        if dirs and root != ACTUAL_DIR:
            return [os.path.join(root, dir) for dir in dirs]

def fix_static_path(static_relative_path, content):
    """Fix the static path in the content.

    Parameters
    ----------
    static_relative_path : str
        The relative path to the static files.
    content : str
        The content to fix the static path in.

    Returns
    -------
    str
        The content with the fixed static path.
    """
    # Replace the static file paths in content href and src attributes having */_static/*
    modified_content = re.sub(
        r'(href|src)="([^:"]*?)/_static/([^"]*?)"',
        r'\1="' + static_relative_path + r'/\3"',
        content
    )
    return modified_content

def fix_static_path_in_file(file_path, static_relative_path):
    """Fix the static path in the file.

    Parameters
    ----------
    file_path : str
        The path of the file.
    static_relative_path : str
        The relative path to the static files.
    """
    with open(file_path, 'r') as f:
        content = f.read()

    modified_content = fix_static_path(static_relative_path, content)

    with open(file_path, 'w') as f:
        f.write(modified_content)

def fix_static_path_in_subfolder_files():
    """Fix the static path in the files."""

    # find the sub folders
    sub_dirs = find_subfolders()
    if not sub_dirs:
        print("No subfolders found")
        return

    # find the sub sub folders
    for sub_dir in sub_dirs:
        sub_dir_files = pathlib.Path(sub_dir).rglob('*.html')
        for file in sub_dir_files:
            fix_static_path_in_file(file, get_relative_path_to_static_files(sub_dir).replace("\\", "/"))

fix_static_path_in_subfolder_files()
