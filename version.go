package gversions

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type suffixKind int

const (
	suffixPrerelease suffixKind = iota
	suffixNone
	suffixPostrelease
)

// Options controls suffix classification and ordering.
type Options struct {
	// PrereleaseSuffixOrder defines prerelease suffix ordering (small -> large).
	// Example: 1.0.0-alpha < 1.0.0-beta < 1.0.0-rc < 1.0.0
	PrereleaseSuffixOrder []string

	// PostreleaseSuffixOrder defines postrelease suffix ordering (small -> large).
	// Example: 1.0.0 < 1.0.0-hotfix < 1.0.0-hotfix2
	PostreleaseSuffixOrder []string
}

var DefaultOptions = Options{
	PrereleaseSuffixOrder: []string{"alpha", "beta", "rc"},
	PostreleaseSuffixOrder: []string{
		"hotfix",
	},
}

// Compare compares two version-like strings.
// Returns -1 if a<b, 0 if a==b, +1 if a>b.
func Compare(a, b string) int {
	return CompareWithOptions(a, b, Options{})
}

// CompareWithOptions compares versions using the provided options.
// If an option slice is nil, the corresponding DefaultOptions slice is used.
func CompareWithOptions(a, b string, opt Options) int {
	prereleaseOrder := opt.PrereleaseSuffixOrder
	postreleaseOrder := opt.PostreleaseSuffixOrder
	if prereleaseOrder == nil {
		prereleaseOrder = DefaultOptions.PrereleaseSuffixOrder
	}
	if postreleaseOrder == nil {
		postreleaseOrder = DefaultOptions.PostreleaseSuffixOrder
	}

	normalize := func(s string) string {
		s = strings.TrimSpace(s)
		// Strip v/V prefix.
		s = strings.TrimPrefix(s, "v")
		s = strings.TrimPrefix(s, "V")
		// Strip product prefix (everything before the first digit).
		if idx := strings.IndexFunc(s, func(r rune) bool { return unicode.IsDigit(r) }); idx > 0 {
			s = s[idx:]
		}
		// Treat underscore as dot.
		s = strings.ReplaceAll(s, "_", ".")
		// '+' handling:
		// - If the version already contains dots, treat '+' as SemVer build metadata delimiter and ignore it.
		// - Otherwise, treat '+' as a legacy separator (e.g. "1+0+0" == "1.0.0").
		if idx := strings.IndexByte(s, '+'); idx >= 0 {
			if strings.Contains(s, ".") {
				s = s[:idx]
			} else {
				s = strings.ReplaceAll(s, "+", ".")
			}
		}
		return s
	}

	a = normalize(a)
	b = normalize(b)

	if a == b {
		return 0
	}

	rankByOrder := func(s string, order []string) int {
		s = strings.ToLower(s)
		for i, token := range order {
			if token == "" {
				continue
			}
			if strings.Contains(s, token) {
				return i + 1
			}
		}
		return 0
	}

	suffixKindOf := func(suffix string) suffixKind {
		suffix = strings.ToLower(suffix)
		if strings.TrimSpace(suffix) == "" {
			return suffixNone
		}
		if rankByOrder(suffix, postreleaseOrder) > 0 {
			return suffixPostrelease
		}
		if rankByOrder(suffix, prereleaseOrder) > 0 {
			return suffixPrerelease
		}
		// Default: unknown suffix sorts after stable (postrelease).
		return suffixPostrelease
	}

	t1 := splitByDotOrDash(a)
	t2 := splitByDotOrDash(b)

	numericPrefixLen := func(tokens []string) int {
		for i := 0; i < len(tokens); i++ {
			if _, err := strconv.Atoi(tokens[i]); err != nil {
				return i
			}
		}
		return len(tokens)
	}

	n1 := numericPrefixLen(t1)
	n2 := numericPrefixLen(t2)

	// Compare numeric prefix (supports 1.0.0.1 style).
	for i := 0; i < n1 && i < n2; i++ {
		aNum, _ := strconv.Atoi(t1[i])
		bNum, _ := strconv.Atoi(t2[i])
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
	}

	// Longer numeric prefix wins: 1.0.0.1 > 1.0.0
	if n1 < n2 {
		return -1
	}
	if n1 > n2 {
		return 1
	}

	suffix1 := strings.Join(t1[n1:], "-")
	suffix2 := strings.Join(t2[n2:], "-")
	kind1 := suffixKindOf(suffix1)
	kind2 := suffixKindOf(suffix2)

	// prerelease < none < postrelease
	if kind1 != kind2 {
		if kind1 < kind2 {
			return -1
		}
		return 1
	}

	if kind1 == suffixPrerelease {
		p1 := suffixPriority(suffix1, prereleaseOrder)
		p2 := suffixPriority(suffix2, prereleaseOrder)
		if p1 > 0 && p2 > 0 && p1 != p2 {
			if p1 < p2 {
				return -1
			}
			return 1
		}
	}

	if kind1 == suffixPostrelease {
		p1 := suffixPriority(suffix1, postreleaseOrder)
		p2 := suffixPriority(suffix2, postreleaseOrder)
		if p1 > 0 && p2 > 0 && p1 != p2 {
			if p1 < p2 {
				return -1
			}
			return 1
		}
	}

	// Compare remaining tokens.
	maxLen := len(t1)
	if len(t2) > maxLen {
		maxLen = len(t2)
	}
	for i := 0; i < maxLen; i++ {
		if i >= len(t1) {
			return -1
		}
		if i >= len(t2) {
			return 1
		}

		ta := t1[i]
		tb := t2[i]
		if ta == tb {
			continue
		}

		aNum, aErr := strconv.Atoi(ta)
		bNum, bErr := strconv.Atoi(tb)
		if aErr == nil && bErr == nil {
			if aNum < bNum {
				return -1
			}
			return 1
		}
		if aErr == nil && bErr != nil {
			return 1
		}
		if aErr != nil && bErr == nil {
			return -1
		}
		if cmp := strings.Compare(ta, tb); cmp != 0 {
			return cmp
		}
	}

	return 0
}

var dotOrDash = regexp.MustCompile(`[\.-]`)

func splitByDotOrDash(s string) []string {
	return dotOrDash.Split(s, -1)
}

func suffixPriority(suffix string, order []string) int {
	suffix = strings.ToLower(suffix)
	for i, token := range order {
		if token == "" {
			continue
		}
		if strings.Contains(suffix, token) {
			return i + 1
		}
	}
	return 0
}
