#!/usr/bin/env python3
"""Keeps CHANGELOG.md in sync with GitHub releases — no manual editing.

Commands:
  sync          Rewrite the [Unreleased] section of CHANGELOG.md from the
                body of the current draft release. No-op if the draft does
                not exist or has no entries (idempotent, so a push-triggered
                commit does not loop).
  finalize      On release publish: move the [Unreleased] content into a
                versioned section "[tag] - YYYY-MM-DD", reset [Unreleased],
                and refresh the compare links at the bottom of the file.

Environment: GITHUB_TOKEN and GITHUB_REPOSITORY (owner/repo).
"""

import os
import re
import subprocess
import sys
from datetime import date

CHANGELOG = "CHANGELOG.md"
CHANGES_MARKER = "## Changes"
FULL_CHANGELOG_MARKER = "**Full Changelog**"
UNRELEASED = "## [Unreleased]"
HEADING = re.compile(r"^## \[", re.MULTILINE)


def sh(*args: str) -> str:
    return subprocess.run(
        args, check=True, capture_output=True, text=True
    ).stdout.strip()


def gh_api(path: str) -> str:
    env = dict(os.environ)
    env["GH_TOKEN"] = os.environ["GITHUB_TOKEN"]
    return subprocess.run(
        ["gh", "api", path], check=True, capture_output=True, text=True, env=env
    ).stdout


def draft_release_body() -> str | None:
    try:
        out = gh_api(f"repos/{os.environ['GITHUB_REPOSITORY']}/releases?per_page=100")
    except subprocess.CalledProcessError:
        return None
    for release in __import__("json").loads(out):
        if release.get("draft"):
            return release.get("body") or ""
    return None


def changes_from_body(body: str) -> str:
    if CHANGES_MARKER not in body:
        return ""
    tail = body.split(CHANGES_MARKER, 1)[1]
    if FULL_CHANGELOG_MARKER in tail:
        tail = tail.split(FULL_CHANGELOG_MARKER, 1)[0]
    return tail.strip()


def read_changelog() -> str:
    with open(CHANGELOG, encoding="utf-8") as f:
        return f.read()


def write_changelog(text: str) -> None:
    with open(CHANGELOG, "w", encoding="utf-8") as f:
        f.write(text)


def find_section(text: str, heading: str) -> tuple[int, int]:
    """Return (start, end) byte offsets of a section whose first line is
    `heading`; end points at the next "## [" line or EOF."""
    start = text.index(heading)
    rest = text[start + len(heading):]
    nxt = HEADING.search(rest)
    end = start + len(heading) + (nxt.start() if nxt else len(rest))
    return start, end


def replace_section(text: str, heading: str, body: str) -> str:
    start, end = find_section(text, heading)
    new = f"{heading}\n\n{body}\n\n"
    return text[:start] + new + text[end:]


def cmd_sync() -> int:
    body = draft_release_body()
    changes = changes_from_body(body) if body else ""
    if not changes:
        print("draft release has no entries — leaving [Unreleased] as-is")
        return 0
    text = read_changelog()
    updated = replace_section(text, UNRELEASED, changes)
    if updated == text:
        print("changelog already up to date")
        return 0
    write_changelog(updated)
    print("updated [Unreleased] from draft release")
    return 0


def previous_tag(tag: str) -> str | None:
    tags = sh("git", "tag", "--sort=-version:refname").splitlines()
    return next((t for t in tags if t != tag), None)


def link_for(tag: str) -> str:
    repo = os.environ["GITHUB_REPOSITORY"]
    prev = previous_tag(tag)
    if prev:
        return f"[{tag}]: https://github.com/{repo}/compare/{prev}...{tag}"
    return f"[{tag}]: https://github.com/{repo}/releases/tag/{tag}"


def cmd_finalize(tag: str) -> int:
    text = read_changelog()
    start, end = find_section(text, UNRELEASED)
    body = text[start + len(UNRELEASED):end].strip()
    if not body:
        print("[Unreleased] is empty — nothing to move")
        return 0

    versioned = f"## [{tag}] - {date.today().isoformat()}\n\n{body}\n\n"
    empty = f"{UNRELEASED}\n\nNo changes yet.\n\n"
    text = text[:start] + empty + versioned + text[end:]

    link = link_for(tag)
    text = re.sub(rf"^\[{re.escape(tag)}\]:.*$", "", text, flags=re.MULTILINE)
    text = re.sub(r"^\[unreleased\]:.*$", f"[unreleased]: https://github.com/{os.environ['GITHUB_REPOSITORY']}/compare/{tag}...HEAD", text, flags=re.MULTILINE)
    text = text.rstrip() + f"\n{link}\n"

    if text.strip() == read_changelog().strip():
        print("changelog already finalized")
        return 0
    write_changelog(text)
    print(f"finalized [{tag}] section")
    return 0


def main() -> int:
    args = sys.argv[1:]
    if not args:
        print(__doc__)
        return 2
    if args[0] == "sync":
        return cmd_sync()
    if args[0] == "finalize":
        tag = None
        for i, a in enumerate(args[1:], 1):
            if a == "--tag" and i + 1 < len(args):
                tag = args[i + 1]
        if not tag:
            print("finalize requires --tag <vX.Y.Z>")
            return 2
        return cmd_finalize(tag)
    print(f"unknown command: {args[0]}")
    return 2


if __name__ == "__main__":
    sys.exit(main())
