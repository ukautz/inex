[![Build Status](https://travis-ci.org/ukautz/inex.svg?branch=master)](https://travis-ci.org/ukautz/inex)
[![Coverage](https://gocover.io/_badge/github.com/ukautz/inex?v=0.1.6)](http://gocover.io/github.com/ukautz/inex)
[![GoDoc](https://godoc.org/github.com/ukautz/inex?status.svg)](https://godoc.org/github.com/ukautz/inex)

Inex : INclude EXclude framework
================================

inex is a Go library to create and execute arbitrary deep boolean filter. Provided is a framework of easy to use primitive operations and types to to build complex and deep, chained filter expressions.

This package is part of my common pattern solutions, which I often (enough) face and hence put in a package at some point.

Dependencies
------------

This package has no dependencies, aside from [testify](https://github.com/stretchr/testify) for testing.

How to use
----------

Install via:

```bash
$ go get github.com/ukautz/inex
```

A few lines of code say more than 1000 words:

```go
package yourpackage

import "github.com/ukautz/inex"

var whitelist = []string{"foo", "bar", "baz", "..."}

func isWhitelisted(str string, excl []string) bool {
	filter := inex.NewRoot().Include(inex.StringsMatcher(whitelist))
	if len(excl) > 0 {
		filter = filter.Exclude(inex.StringsMatcher(excl))
	}
	return filter.Root().Match(str)
}
```

Matcher
-------

The `Matcher` interface is the core of this package. It requires the `Match(string) bool` method to be implemented.

Provided are a list of essential implementations like `StringMatcher` (matching a single string), `StringsMatcher` (matching a slice of strings), `RegexpMatcher` (matching a regular expression) and a `FuncMatcher` (matching `func(string) bool`).

This package comes also with a host of tooling functions like `And(...Matcher)` (matches if all Matchers match), `Or(...Matcher)` (matches if any Matcher matches), `Not(Matcher)` (inverts matcher) and so on.

Example: Rebuild GNU linux `find` - or parts of it
--------------------------------------------------

Let's build a simple command line tool to recursively find all files in a directory, which are included or excluded by provided regular expressions.

The command line interface should provide `-i <regex>` or `--include <regex>` parameters to include matching paths and `-e <regex>` or `--exclude <regex>` to do the opposite. The exemplary `filter` command line program, executed in this `inex` directory should then work like this:

```bash
$ newfind -i '.*\.go$' -i '(_test|inex)' -e '\/'
inex.go
inex_test.go
matcher_test.go
```

This could be done in a few lines:

```go
package main

import (
	"fmt"
	"github.com/ukautz/inex"
	"os"
	"path/filepath"
	"regexp"
)

var help = fmt.Sprintf("Usage: %s ((-e|-i|--exclude|--include) <regex>)+\n", os.Args[0])

func main() {
	l := len(os.Args)
	if l <= 1 || l%2 != 1 {
		panic(help)
	}
	filter := inex.NewRoot()
	for i := 1; i < l; i += 2 {
		matcher := &inex.RegexpMatcher{regexp.MustCompile(os.Args[i+1])}
		switch os.Args[i] {
		case "-e", "--exclude":
			filter = filter.Exclude(matcher)
		case "-i", "--include":
			filter = filter.Include(matcher)
		default:
			panic(fmt.Sprintf("unsupported operation \"%s\"\n%s", os.Args[i], help))
		}
	}
	filter = filter.Root()
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if filter.Match(path) {
			fmt.Println(path)
		}
		return nil
	})
}
```

See even more examples in in `example/newfinder.go`
