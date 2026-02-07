# gversions

[![Go Reference](https://pkg.go.dev/badge/github.com/aak1247/gversions.svg)](https://pkg.go.dev/github.com/aak1247/gversions)
[![Go Report Card](https://goreportcard.com/badge/github.com/aak1247/gversions)](https://goreportcard.com/report/github.com/aak1247/gversions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/aak1247/gversions)](go.mod)
[![Release](https://img.shields.io/github/v/release/aak1247/gversions?sort=semver)](https://github.com/aak1247/gversions/releases)

[English](README.md) | [简体中文](README.zh-CN.md)

`github.com/aak1247/gversions` 是一个**版本号比较器**：数值核心（`MAJOR.MINOR.PATCH`）遵循 SemVer 的比较逻辑，同时补充了在 Git tag 工作流中经常用到的概念：**postrelease（发布后后缀）**，例如 `-hotfix`，它会排在稳定版之后。

## 安装

```bash
go get github.com/aak1247/gversions@latest
```

```go
import "github.com/aak1247/gversions"
```

## 快速开始

```go
gversions.Compare("1.2.3-hotfix1", "1.2.3") // => 1
gversions.Compare("1.2.3-rc.1", "1.2.3")    // => -1
```

对 tag 排序：

```go
import (
	"sort"

	"github.com/aak1247/gversions"
)

sort.Slice(tags, func(i, j int) bool {
	return gversions.Compare(tags[i], tags[j]) < 0
})
```

## 为什么需要 postrelease？

SemVer 定义了 **prerelease（预发布）** 标识（例如 `-alpha`、`-beta`、`-rc`），其排序规则是：**稳定版之前**：

- `1.2.3-alpha` < `1.2.3`
- `1.2.3-rc.1` < `1.2.3`

但在实际的发布 / 打 tag 过程中，团队也常用一些后缀表达“在不提升 `MAJOR.MINOR.PATCH` 的情况下，比稳定版更新一点”的含义，例如：

- `1.2.3-hotfix` > `1.2.3`
- `1.2.3-hotfix2` > `1.2.3-hotfix1`

SemVer 本身并没有“发布后(postrelease)”的比较规则（build metadata 的 `+...` 明确不参与排序），因此这里补齐了这类 tag 的比较与选择需求。

## 排序模型

对于同一个数值核心（例如 `1.2.3`），后缀会被划分成三类：

1. **prerelease**：排在稳定版之前（`alpha < beta < rc < stable`）
2. **stable**：无后缀
3. **postrelease**：排在稳定版之后（`stable < hotfix < hotfix2 ...`）

默认规则：

- prerelease tokens：`alpha`、`beta`、`rc`
- postrelease tokens：`hotfix`
- 未知后缀：默认按 **postrelease** 处理（也就是排在稳定版之后）

你可以通过 `Options` 覆盖这些排序规则。

## 自定义排序

```go
opt := gversions.Options{
	PrereleaseSuffixOrder:  []string{"rc", "beta", "alpha"}, // rc < beta < alpha
	PostreleaseSuffixOrder: []string{"p1", "p2"},
}
gversions.CompareWithOptions("1.0.0-p2", "1.0.0-p1", opt) // => 1
```

## 规范化（更适合 tag）

这个比较器面向 tag 字符串，不要求完全严格的 SemVer：

- 忽略前导 `v` / `V`（`v1.2.3` == `1.2.3`）
- 忽略 build metadata 的 `+...`（`1.2.3+build.1` == `1.2.3+build.2`）
- 对于没有 `.` 的旧式 tag，把 `+` 当作分隔符（`1+2+3` == `1.2.3`）
- 如果第一个数字之前存在产品前缀，会被忽略（`hive.1.2.3` == `1.2.3`）
- `_` 会当作 `.`（`1_2_3` == `1.2.3`）

## 严格 SemVer

如果你需要严格的 SemVer 排序，请使用 `CompareSemver`（基于 `golang.org/x/mod/semver` 的轻量封装）。
