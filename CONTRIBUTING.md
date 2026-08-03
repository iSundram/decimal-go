# Contributing to decimal-go

Thank you for your interest in contributing to decimal-go! We welcome all contributions, from bug fixes to new features and documentation improvements.

## How to Contribute

1. **Fork the repository**: Click the "Fork" button at the top right of the repository page.
2. **Clone your fork**: `git clone https://github.com/<your-username>/decimal-go.git`
3. **Create a branch**: `git checkout -b my-new-feature`
4. **Make your changes**: Write your code and add tests if applicable.
5. **Run tests**: Ensure all tests pass by running `go test ./...`.
6. **Commit your changes**: `git commit -am 'Add some feature'`
7. **Push to the branch**: `git push origin my-new-feature`
8. **Submit a pull request**: Open a pull request against the `main` branch of the original repository.

## Coding Style

- Please adhere to standard Go formatting. Run `go fmt` before committing.
- Ensure any new logic is well-documented with GoDoc comments.
- Since this library aims to mirror `decimal.js`, please ensure new features or bug fixes maintain parity with its behavior.

## Reporting Issues

If you find a bug or have a feature request, please open an issue on GitHub. Include as much detail as possible, such as code snippets and expected vs. actual behavior.

## Releases & Changelog

Releases and the changelog are fully automated — no manual editing:

1. **Label your PR** (so it is grouped correctly): `feature`, `enhancement`,
   `fix`, `bug`, `documentation`, `docs`, `tests`, `quality`, or `ci`.
   Unlabelled PRs default to the `patch` group. Prefix labels like `docs:`
   are applied automatically by the autolabeler.
2. **Merge the PR.** The [Release Drafter](.github/workflows/release-drafter.yml)
   workflow updates the draft release notes, and the
   [Update CHANGELOG](.github/workflows/update-changelog.yml) workflow
   rewrites the `[Unreleased]` section of `CHANGELOG.md` automatically.
3. **Publish a release**: open a PR titled `release` with the `release`
   label. When it merges, the draft is published as a versioned release
   (version bumped from the labels), and the
   [Finalize CHANGELOG](.github/workflows/finalize-changelog.yml) workflow
   moves `[Unreleased]` into the versioned section with the release date.

Add a `skip-changelog` label to any PR you do not want mentioned.
