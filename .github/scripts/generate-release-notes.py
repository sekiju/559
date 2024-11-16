import pathlib
import re

IN = "CHANGELOG.md"
OUT = "whats-new.md"

script_path = pathlib.Path(__file__).resolve()
project_root = script_path.parent.parent.parent

with open(pathlib.Path(project_root, IN), "r") as f:
    changelog_content = f.read()

    version_sections = re.split(r'(?=## v\d+)', changelog_content)
    version_sections = [s for s in version_sections if not ('<!--' in s or '-->' in s)]
    if not version_sections:
        raise ValueError("No valid version sections found in changelog")

    lines = version_sections[0].split('\n')

    filtered_lines = []
    for line in lines:
        if line.startswith('## v'):
            continue
        if line.startswith('Date:'):
            continue
        if line.strip() or filtered_lines:
            filtered_lines.append(line)

    current_changelog = '\n'.join(filtered_lines).strip()

    with open(pathlib.Path(project_root, OUT), "w") as ft:
        ft.write(current_changelog)