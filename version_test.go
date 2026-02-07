package gversions

import "testing"

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected int
	}{
		{name: "Equal", a: "1.0.0", b: "1.0.0", expected: 0},
		{name: "Major", a: "2.0.0", b: "1.0.0", expected: 1},
		{name: "Minor", a: "1.2.0", b: "1.1.0", expected: 1},
		{name: "Patch", a: "1.0.2", b: "1.0.1", expected: 1},
		{name: "V prefix ignored", a: "v1.0.0", b: "1.0.0", expected: 0},
		{name: "Product prefix ignored", a: "hive.1.0.0", b: "1.0.0", expected: 0},
		{name: "Underscore treated as dot", a: "1_0_0", b: "1.0.0", expected: 0},
		{name: "Legacy plus treated as dot", a: "1+0+0", b: "1.0.0", expected: 0},
		{name: "Build metadata ignored", a: "1.0.0+build.1", b: "1.0.0+build.2", expected: 0},
		{name: "Build metadata ignored vs stable", a: "1.0.0+build.1", b: "1.0.0", expected: 0},
		{name: "Longer numeric wins", a: "1.0.0.1", b: "1.0.0", expected: 1},
		{name: "Prerelease less than stable", a: "1.0.0-alpha", b: "1.0.0", expected: -1},
		{name: "RC greater than beta", a: "1.0.0-rc.1", b: "1.0.0-beta.2", expected: 1},
		{name: "Postrelease greater than stable", a: "1.0.0-hotfix", b: "1.0.0", expected: 1},
		{name: "Unknown suffix defaults postrelease", a: "1.0.0-foo", b: "1.0.0", expected: 1},
		{name: "Suffix numeric compare", a: "1.0.0-hotfix2", b: "1.0.0-hotfix1", expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compare(tt.a, tt.b)
			if got != tt.expected {
				t.Fatalf("Compare(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestCompareSymmetry(t *testing.T) {
	versions := []string{
		"1.0.0",
		"v1.0.0",
		"hive.1.0.0",
		"1.0.1",
		"1.0.2-alpha",
		"1.0.2-beta",
		"1.0.2-rc",
		"1.0.2",
		"1.0.2-hotfix1",
		"1.1.0",
		"2.0.0",
	}

	for _, a := range versions {
		for _, b := range versions {
			ab := Compare(a, b)
			ba := Compare(b, a)
			if ab != -ba {
				t.Fatalf("symmetry failed: Compare(%q,%q)=%d Compare(%q,%q)=%d", a, b, ab, b, a, ba)
			}
		}
	}
}

func TestCompareWithOptions_CustomPrereleaseOrder(t *testing.T) {
	// Reverse the default prerelease order: rc < beta < alpha
	opt := Options{
		PrereleaseSuffixOrder: []string{"rc", "beta", "alpha"},
	}
	if CompareWithOptions("1.0.0-alpha", "1.0.0-rc", opt) <= 0 {
		t.Fatalf("expected alpha > rc with custom prerelease order")
	}
}

func TestCompareWithOptions_CustomPostreleaseOrder(t *testing.T) {
	opt := Options{
		PostreleaseSuffixOrder: []string{"p1", "p2"},
	}
	if CompareWithOptions("1.0.0-p2", "1.0.0-p1", opt) <= 0 {
		t.Fatalf("expected p2 > p1 with custom postrelease order")
	}
}

func BenchmarkCompare(b *testing.B) {
	a := "hive.2.1.3-beta.2"
	c := "hive.2.1.3-beta.1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Compare(a, c)
	}
}

func TestCompareSemver(t *testing.T) {
	if got := CompareSemver("1.0.0-foo", "1.0.0"); got >= 0 {
		t.Fatalf("CompareSemver prerelease should be < stable, got %d", got)
	}
	if got := CompareSemver("1.0.0+build.1", "1.0.0+build.2"); got != 0 {
		t.Fatalf("CompareSemver build metadata ignored, got %d", got)
	}
}
