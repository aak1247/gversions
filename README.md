# gverions

`github.com/aak1247/gverions` provides a version comparator that is **SemVer-aligned** for the numeric core, and extends it with a practical concept often needed in Git tag workflows: **postrelease** suffixes.

## Why postrelease?

SemVer defines **prerelease** identifiers (e.g. `-alpha`, `-beta`, `-rc`) which sort **before** the stable release:

- `1.2.3-alpha` < `1.2.3`
- `1.2.3-rc.1` < `1.2.3`

In real-world release/tagging practices, teams sometimes need a suffix that sorts **after** the stable release without bumping `MAJOR.MINOR.PATCH`, for example:

- `1.2.3-hotfix` > `1.2.3`
- `1.2.3-hotfix2` > `1.2.3-hotfix1`

SemVer itself does not provide a “postrelease” precedence rule (build metadata `+...` is explicitly ignored for ordering), so this package fills that gap for tag ordering and selection.

## Ordering model

For a given numeric core (e.g. `1.2.3`), suffixes are classified into three kinds:

1. **prerelease**: sorts before stable (`alpha < beta < rc < stable`)
2. **stable**: no suffix
3. **postrelease**: sorts after stable (`stable < hotfix < hotfix2 ...`)

By default:

- prerelease tokens: `alpha`, `beta`, `rc`
- postrelease tokens: `hotfix`
- unknown suffixes: treated as **postrelease** (so they sort after stable)

You can override ordering via `Options`.

## Normalization

The comparator is designed for tag strings, not only strict SemVer:

- Leading `v` / `V` is ignored (`v1.2.3` == `1.2.3`)
- SemVer build metadata `+...` is ignored for ordering (`1.2.3+build.1` == `1.2.3+build.2`)
- Product prefix before the first digit is ignored (`hive.1.2.3` == `1.2.3`)
- `_` is treated as `.` (`1_2_3` == `1.2.3`)

## API

```go
import "github.com/aak1247/gverions"

gverions.Compare("1.2.3-hotfix1", "1.2.3") // => 1
gverions.Compare("1.2.3-rc.1", "1.2.3")    // => -1
```

Custom ordering:

```go
opt := gverions.Options{
  PrereleaseSuffixOrder: []string{"rc", "beta", "alpha"}, // rc < beta < alpha
  PostreleaseSuffixOrder: []string{"p1", "p2"},
}
gverions.CompareWithOptions("1.0.0-p2", "1.0.0-p1", opt) // => 1
```

## Notes / Differences from strict SemVer

This package intentionally deviates from strict SemVer in one place: **a hyphen suffix can be treated as postrelease** if it matches `PostreleaseSuffixOrder` (or is unknown, by default).

If you need strict SemVer precedence, use `golang.org/x/mod/semver` directly.
