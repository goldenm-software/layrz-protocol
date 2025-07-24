"""Deploy"""

import os
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent
VERSION = os.environ.get("VERSION", "v0.0.0").replace("v", "")
toml_path = BASE_DIR / "pubspec.yaml"

with open(toml_path, "r") as f:
    lines = f.readlines()

line_no = None
for i, line in enumerate(lines):
    if "version" in line:
        line_no = i
        break

if line_no is None:
    raise ValueError("Version not found in pyproject.toml")

lines[line_no] = f'version: "{VERSION.strip()}"\n'

with open(toml_path, "w") as f:
    f.writelines(lines)
