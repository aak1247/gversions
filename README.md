# gversions

[![Go Reference](https://pkg.go.dev/badge/github.com/aak1247/gversions.svg)](https://pkg.go.dev/github.com/aak1247/gversions)
[![Go Report Card](https://goreportcard.com/badge/github.com/aak1247/gversions)](https://goreportcard.com/report/github.com/aak1247/gversions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/aak1247/gversions)](go.mod)
[![Release](https://img.shields.io/github/v/release/aak1247/gversions?sort=semver)](https://github.com/aak1247/gversions/releases)

[English](README.md) | [简体中文](README.zh-CN.md)

`github.com/aak1247/gversions` is a version comparator that is **SemVer-aligned** on the numeric core (`MAJOR.MINOR.PATCH`), and adds a practical concept often needed in Git tag workflows: **postrelease** suffixes (e.g. `-hotfix`), which sort **after** the stable release.

## Install

```bash
go get github.com/aak1247/gversions@latest
```

```go
import "github.com/aak1247/gversions"
```

## Quick start

```go
gversions.Compare("1.2.3-hotfix1", "1.2.3") // => 1
gversions.Compare("1.2.3-rc.1", "1.2.3")    // => -1
```

Sort tags:

```go
import (
	"sort"

	"github.com/aak1247/gversions"
)

sort.Slice(tags, func(i, j int) bool {
	return gversions.Compare(tags[i], tags[j]) < 0
})
```

## Why “postrelease”?

SemVer defines **prerelease** identifiers (e.g. `-alpha`, `-beta`, `-rc`) which sort **before** the stable release:

- `1.2.3-alpha` < `1.2.3`
- `1.2.3-rc.1` < `1.2.3`

In many real-world tagging/release practices, teams also use suffixes that should sort **after** the stable release without bumping `MAJOR.MINOR.PATCH`, for example:

- `1.2.3-hotfix` > `1.2.3`
- `1.2.3-hotfix2` > `1.2.3-hotfix1`

SemVer does not define a “postrelease” precedence rule (build metadata `+...` is ignored for ordering), so this package fills that gap for tag ordering and selection.

## Ordering model

For a given numeric core (e.g. `1.2.3`), suffixes are classified into three kinds:

1. **prerelease**: sorts before stable (`alpha < beta < rc < stable`)
2. **stable**: no suffix
3. **postrelease**: sorts after stable (`stable < hotfix < hotfix2 ...`)

Defaults:

- prerelease tokens: `alpha`, `beta`, `rc`
- postrelease tokens: `hotfix`
- unknown suffixes: treated as **postrelease** (so they sort after stable)

You can override ordering via `Options`.

## Custom ordering

```go
opt := gversions.Options{
	PrereleaseSuffixOrder:  []string{"rc", "beta", "alpha"}, // rc < beta < alpha
	PostreleaseSuffixOrder: []string{"p1", "p2"},
}
gversions.CompareWithOptions("1.0.0-p2", "1.0.0-p1", opt) // => 1
```

## Normalization (tag-friendly)

This comparator is designed for tag strings, not only strict SemVer:

- Leading `v` / `V` is ignored (`v1.2.3` == `1.2.3`)
- SemVer build metadata `+...` is ignored for ordering (`1.2.3+build.1` == `1.2.3+build.2`)
- For legacy tags without dots, `+` is treated as a separator (`1+2+3` == `1.2.3`)
- Product prefix before the first digit is ignored (`hive.1.2.3` == `1.2.3`)
- `_` is treated as `.` (`1_2_3` == `1.2.3`)

## Strict SemVer

If you need strict SemVer precedence, use `CompareSemver` (a thin wrapper around `golang.org/x/mod/semver`).
