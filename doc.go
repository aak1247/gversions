// Package gversions provides semver-aligned version comparison with support for
// "postrelease" suffixes (e.g. hotfix).
//
// SemVer defines prerelease suffixes like "-alpha", "-beta", "-rc" which sort
// *before* the stable release. In many Git tag workflows, you may also need a
// suffix that sorts *after* the stable release without bumping MAJOR/MINOR/PATCH,
// such as "1.2.3-hotfix1" > "1.2.3".
//
// This package keeps semver-like numeric comparison for the core version, and
// lets you classify suffix tokens into:
//   - prerelease: sorts before the stable release (alpha < beta < rc < stable)
//   - postrelease: sorts after the stable release (stable < hotfix < hotfix2 ...)
//
// Notes:
//   - Leading "v"/"V" is ignored (v1.2.3 == 1.2.3).
//   - If a tag has a product prefix before the first digit (e.g. "hive.1.2.3"),
//     the prefix is ignored for comparison.
//   - Unknown suffixes are treated as postrelease by default (so they sort after
//     stable) to match common "tag" expectations. You can override ordering via Options.
package gversions
