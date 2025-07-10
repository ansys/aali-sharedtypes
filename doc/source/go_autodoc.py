#!/usr/bin/env python3

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

"""
Go AutoDoc Generator for Sphinx

This script parses Go source files and generates Sphinx autodoc-compatible
documentation for the AALI SharedTypes API reference.
"""

import os
import re
import glob
from pathlib import Path
from typing import Dict, List, Tuple, Optional

class GoParser:
    """Parser for Go source files to extract function and type information."""

    def __init__(self, source_dir: str):
        self.source_dir = Path(source_dir)
        self.packages = {}

    def parse_package(self, package_path: str) -> Dict:
        """Parse a Go package and extract all functions, types, and interfaces."""
        package_dir = self.source_dir / package_path
        if not package_dir.exists():
            return {}

        package_info = {
            'name': package_path.replace('/', '.'),
            'functions': [],
            'types': [],
            'interfaces': [],
            'constants': [],
            'variables': []
        }

        # Find all .go files in the package
        go_files = list(package_dir.glob('*.go'))

        for go_file in go_files:
            self._parse_go_file(go_file, package_info)

        return package_info

    def _parse_go_file(self, file_path: Path, package_info: Dict):
        """Parse a single Go file and extract its contents."""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
        except Exception as e:
            print(f"Warning: Could not read {file_path}: {e}")
            return

        # Extract package name
        package_match = re.search(r'^package\s+(\w+)', content, re.MULTILINE)
        if package_match:
            package_info['package_name'] = package_match.group(1)

        # Extract functions
        functions = self._extract_functions(content)
        package_info['functions'].extend(functions)

        # Extract types
        types = self._extract_types(content)
        package_info['types'].extend(types)

        # Extract interfaces
        interfaces = self._extract_interfaces(content)
        package_info['interfaces'].extend(interfaces)

        # Extract constants
        constants = self._extract_constants(content)
        package_info['constants'].extend(constants)

        # Extract variables
        variables = self._extract_variables(content)
        package_info['variables'].extend(variables)

    def _extract_functions(self, content: str) -> List[Dict]:
        """Extract function definitions from Go code."""
        functions = []

        # Match function definitions
        # This regex matches both exported and unexported functions
        pattern = r'^func\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(([^)]*)\)\s*([^{]*?)\s*{'

        for match in re.finditer(pattern, content, re.MULTILINE):
            func_name = match.group(1)
            params = match.group(2).strip()
            return_type = match.group(3).strip()

            # Extract comments above the function
            comments = self._extract_comments_before(content, match.start())

            functions.append({
                'name': func_name,
                'params': params,
                'return_type': return_type,
                'comments': comments,
                'exported': func_name[0].isupper() if func_name else False
            })

        return functions

    def _extract_types(self, content: str) -> List[Dict]:
        """Extract type definitions from Go code."""
        types = []

        # Match type definitions
        pattern = r'^type\s+([A-Za-z_][A-Za-z0-9_]*)\s+([^{]+?)(?:\s*{([^}]*)}|$)'

        for match in re.finditer(pattern, content, re.MULTILINE):
            type_name = match.group(1)
            type_def = match.group(2).strip()
            type_body = match.group(3) if match.group(3) else ""

            comments = self._extract_comments_before(content, match.start())

            types.append({
                'name': type_name,
                'definition': type_def,
                'body': type_body,
                'comments': comments,
                'exported': type_name[0].isupper() if type_name else False
            })

        return types

    def _extract_interfaces(self, content: str) -> List[Dict]:
        """Extract interface definitions from Go code."""
        interfaces = []

        # Match interface definitions
        pattern = r'^type\s+([A-Za-z_][A-Za-z0-9_]*)\s+interface\s*{([^}]*)}'

        for match in re.finditer(pattern, content, re.MULTILINE):
            interface_name = match.group(1)
            interface_body = match.group(2).strip()

            comments = self._extract_comments_before(content, match.start())

            # Parse interface methods
            methods = []
            for line in interface_body.split('\n'):
                line = line.strip()
                if line and not line.startswith('//'):
                    method_match = re.match(r'([A-Za-z_][A-Za-z0-9_]*)\s*\(([^)]*)\)\s*([^{]*)', line)
                    if method_match:
                        methods.append({
                            'name': method_match.group(1),
                            'params': method_match.group(2),
                            'return_type': method_match.group(3).strip()
                        })

            interfaces.append({
                'name': interface_name,
                'methods': methods,
                'comments': comments,
                'exported': interface_name[0].isupper() if interface_name else False
            })

        return interfaces

    def _extract_constants(self, content: str) -> List[Dict]:
        """Extract constant definitions from Go code."""
        constants = []

        # Match const blocks
        const_pattern = r'^const\s*\(([^)]*)\)'

        for match in re.finditer(const_pattern, content, re.MULTILINE):
            const_block = match.group(1)
            comments = self._extract_comments_before(content, match.start())

            # Parse individual constants
            for line in const_block.split('\n'):
                line = line.strip()
                if line and not line.startswith('//'):
                    const_match = re.match(r'([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.+)', line)
                    if const_match:
                        constants.append({
                            'name': const_match.group(1),
                            'value': const_match.group(2).strip(),
                            'comments': comments,
                            'exported': const_match.group(1)[0].isupper() if const_match.group(1) else False
                        })

        return constants

    def _extract_variables(self, content: str) -> List[Dict]:
        """Extract variable definitions from Go code."""
        variables = []

        # Match var blocks
        var_pattern = r'^var\s+([A-Za-z_][A-Za-z0-9_]*)\s+([^=]+?)(?:\s*=\s*(.+))?$'

        for match in re.finditer(var_pattern, content, re.MULTILINE):
            var_name = match.group(1)
            var_type = match.group(2).strip()
            var_value = match.group(3).strip() if match.group(3) else ""

            comments = self._extract_comments_before(content, match.start())

            variables.append({
                'name': var_name,
                'type': var_type,
                'value': var_value,
                'comments': comments,
                'exported': var_name[0].isupper() if var_name else False
            })

        return variables

    def _extract_comments_before(self, content: str, position: int) -> str:
        """Extract comments that appear before a given position."""
        # Look for comments in the lines before the position
        lines = content[:position].split('\n')
        comments = []

        for line in reversed(lines):
            line = line.strip()
            if line.startswith('//'):
                comments.insert(0, line[2:].strip())
            elif line.startswith('/*'):
                # Handle block comments
                pass
            elif line and not line.startswith('//'):
                break

        return ' '.join(comments)

class SphinxDocGenerator:
    """Generate Sphinx documentation from parsed Go code."""

    def __init__(self, output_dir: str):
        self.output_dir = Path(output_dir)
        self.output_dir.mkdir(parents=True, exist_ok=True)

    def generate_package_doc(self, package_info: Dict) -> str:
        """Generate Sphinx documentation for a Go package."""
        if not package_info:
            return ""

        doc_lines = []

        # Package header
        package_name = package_info.get('name', 'Unknown Package')
        doc_lines.append(f".. _{package_name.replace('.', '_')}:")
        doc_lines.append("")
        doc_lines.append(f"{package_name}")
        doc_lines.append("=" * len(package_name))
        doc_lines.append("")

        # Package description
        doc_lines.append(f"This package provides functionality for {package_name}.")
        doc_lines.append("")

        # Functions section
        if package_info.get('functions'):
            doc_lines.append("Functions")
            doc_lines.append("---------")
            doc_lines.append("")

            for func in package_info['functions']:
                if func.get('exported', False):  # Only document exported functions
                    doc_lines.extend(self._generate_function_doc(func))
                    doc_lines.append("")

        # Types section
        if package_info.get('types'):
            doc_lines.append("Types")
            doc_lines.append("-----")
            doc_lines.append("")

            for type_info in package_info['types']:
                if type_info.get('exported', False):  # Only document exported types
                    doc_lines.extend(self._generate_type_doc(type_info))
                    doc_lines.append("")

        # Interfaces section
        if package_info.get('interfaces'):
            doc_lines.append("Interfaces")
            doc_lines.append("----------")
            doc_lines.append("")

            for interface in package_info['interfaces']:
                if interface.get('exported', False):  # Only document exported interfaces
                    doc_lines.extend(self._generate_interface_doc(interface))
                    doc_lines.append("")

        # Constants section
        if package_info.get('constants'):
            doc_lines.append("Constants")
            doc_lines.append("---------")
            doc_lines.append("")

            for const in package_info['constants']:
                if const.get('exported', False):  # Only document exported constants
                    doc_lines.extend(self._generate_constant_doc(const))
                    doc_lines.append("")

        # Variables section
        if package_info.get('variables'):
            doc_lines.append("Variables")
            doc_lines.append("---------")
            doc_lines.append("")

            for var in package_info['variables']:
                if var.get('exported', False):  # Only document exported variables
                    doc_lines.extend(self._generate_variable_doc(var))
                    doc_lines.append("")

        return '\n'.join(doc_lines)

    def _generate_function_doc(self, func: Dict) -> List[str]:
        """Generate documentation for a function."""
        lines = []

        # Function signature - escape the signature to avoid formatting issues
        signature = f"func {func['name']}({func['params']}) {func['return_type']}"
        signature = self._clean_comments(signature)  # Escape the function signature
        lines.append(f".. function:: {func['name']}")
        lines.append("")
        lines.append(f"   {signature}")
        lines.append("")

        # Function description
        if func.get('comments'):
            # Clean up comments to avoid formatting issues
            comments = self._clean_comments(func['comments'])
            lines.append(f"   {comments}")
            lines.append("")

        return lines

    def _generate_type_doc(self, type_info: Dict) -> List[str]:
        """Generate documentation for a type."""
        lines = []

        # Type definition - escape the definition to avoid formatting issues
        type_def = f"type {type_info['name']} {type_info['definition']}"
        type_def = self._clean_comments(type_def)  # Escape the type definition
        lines.append(f".. type:: {type_info['name']}")
        lines.append("")
        lines.append(f"   {type_def}")
        lines.append("")

        # Type description
        if type_info.get('comments'):
            comments = self._clean_comments(type_info['comments'])
            lines.append(f"   {comments}")
            lines.append("")

        return lines

    def _generate_interface_doc(self, interface: Dict) -> List[str]:
        """Generate documentation for an interface."""
        lines = []

        # Interface definition
        lines.append(f".. type:: {interface['name']}")
        lines.append("")

        # Interface description
        if interface.get('comments'):
            comments = self._clean_comments(interface['comments'])
            lines.append(f"   {comments}")
            lines.append("")

        # Interface methods
        if interface.get('methods'):
            lines.append("   **Methods:**")
            lines.append("")
            for method in interface['methods']:
                method_sig = f"{method['name']}({method['params']}) {method['return_type']}"
                method_sig = self._clean_comments(method_sig)  # Escape the method signature
                lines.append(f"   * {method_sig}")
            lines.append("")

        return lines

    def _generate_constant_doc(self, const: Dict) -> List[str]:
        """Generate documentation for a constant."""
        lines = []

        # Constant definition
        const_def = f"const {const['name']} = {const['value']}"
        lines.append(f".. data:: {const['name']}")
        lines.append("")
        lines.append(f"   {const_def}")
        lines.append("")

        # Constant description
        if const.get('comments'):
            comments = self._clean_comments(const['comments'])
            lines.append(f"   {comments}")
            lines.append("")

        return lines

    def _generate_variable_doc(self, var: Dict) -> List[str]:
        """Generate documentation for a variable."""
        lines = []

        # Variable definition - escape the definition to avoid formatting issues
        var_def = f"var {var['name']} {var['type']}"
        if var.get('value'):
            var_def += f" = {var['value']}"

        var_def = self._clean_comments(var_def)  # Escape the variable definition
        lines.append(f".. data:: {var['name']}")
        lines.append("")
        lines.append(f"   {var_def}")
        lines.append("")

        # Variable description
        if var.get('comments'):
            comments = self._clean_comments(var['comments'])
            lines.append(f"   {comments}")
            lines.append("")

        return lines

    def _clean_comments(self, comments: str) -> str:
        """Clean up comments to avoid Sphinx formatting issues."""
        if not comments:
            return ""

        # Replace problematic characters and patterns
        cleaned = comments

        # Escape all asterisks (simpler approach)
        cleaned = cleaned.replace('*', r'\*')

        # Escape all underscores (simpler approach)
        cleaned = cleaned.replace('_', r'\_')

        # Escape backticks
        cleaned = cleaned.replace('`', r'\`')

        # Escape square brackets
        cleaned = cleaned.replace('[', r'\[').replace(']', r'\]')

        # Escape curly braces
        cleaned = cleaned.replace('{', r'\{').replace('}', r'\}')

        # Escape pipe characters
        cleaned = cleaned.replace('|', r'\|')

        # Escape plus signs
        cleaned = cleaned.replace('+', r'\+')

        # Escape equals signs
        cleaned = cleaned.replace('=', r'\=')

        # Escape hash signs
        cleaned = cleaned.replace('#', r'\#')

        # Escape tilde
        cleaned = cleaned.replace('~', r'\~')

        # Escape caret
        cleaned = cleaned.replace('^', r'\^')

        return cleaned

def main():
    """Main function to generate Go documentation."""
    # Configuration
    source_dir = "../../pkg"  # Correct relative path from doc/source/
    output_dir = "advanced/autodoc"

    # Create parser and generator
    parser = GoParser(source_dir)
    generator = SphinxDocGenerator(output_dir)

    # Define packages to document (aali-sharedtypes)
    packages = [
        "aali_graphdb",
        "aaliagentgrpc",
        "aaliflowkitgrpc",
        "clients/flowkitclient",
        "clients/flowkitpythonclient",
        "config",
        "logging",
        "sharedtypes",
        "typeconverters"
    ]

    # Generate documentation for each package
    for package_path in packages:
        print(f"Parsing package: {package_path}")
        package_info = parser.parse_package(package_path)

        if package_info:
            doc_content = generator.generate_package_doc(package_info)

            # Write to file
            output_file = generator.output_dir / f"{package_path.replace('/', '_')}.rst"
            output_file.parent.mkdir(parents=True, exist_ok=True)

            with open(output_file, 'w', encoding='utf-8') as f:
                f.write(doc_content)

            print(f"Generated: {output_file}")
        else:
            print(f"No content found for package: {package_path}")

    # Generate index file
    index_content = generate_index_file(packages)
    index_file = generator.output_dir / "index.rst"

    with open(index_file, 'w', encoding='utf-8') as f:
        f.write(index_content)

    print(f"Generated index: {index_file}")

def generate_index_file(packages: List[str]) -> str:
    """Generate an index file for the autodoc directory."""
    lines = [
        ".. _autodoc_index:",
        "",
        "Auto-Generated API Documentation",
        "=================================",
        "",
        "This section contains auto-generated documentation for all Go packages in AALI SharedTypes.",
        "",
        ".. toctree::",
        "   :maxdepth: 2",
        "   :caption: Packages",
        ""
    ]

    for package in packages:
        file_name = package.replace('/', '_')
        lines.append(f"   {file_name}")

    return '\n'.join(lines)

if __name__ == "__main__":
    main()
