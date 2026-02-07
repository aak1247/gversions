package gversions

import "golang.org/x/mod/semver"

// CanonicalSemver normalizes a SemVer-like string to its canonical form.
//
// - Empty string becomes "v0.0.0"
// - Missing "v" prefix is added
// - Invalid semver is returned as-is
func CanonicalSemver(input string) string {
	if input == "" {
		return "v0.0.0"
	}
	v := input
	if v[0] != 'v' {
		v = "v" + v
	}
	if !semver.IsValid(v) {
		return v
	}
	return semver.Canonical(v)
}

// CompareSemver compares two versions using strict SemVer precedence.
//
// If both inputs are valid semver after CanonicalSemver, semver.Compare is used.
// Otherwise it falls back to string compare to keep deterministic ordering.
func CompareSemver(a, b string) int {
	na := CanonicalSemver(a)
	nb := CanonicalSemver(b)
	if semver.IsValid(na) && semver.IsValid(nb) {
		return semver.Compare(na, nb)
	}
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
