package main

import (
	"regexp"
	"strings"
)

func findLineIdx(content []string, path ...string) int {
	//TODO - content cannot be null, path at least 1 element
	re := make([]*regexp.Regexp, len(path))
	for i, v := range path {
		//TODO - use Compile and handle error
		re[i] = regexp.MustCompile(v)
	}

	reIdx := 0
	for lineIdx, line := range content {
		if re[reIdx].MatchString(line) {
			reIdx++
			if reIdx == len(path) {
				return lineIdx
			}
		}
	}
	return -1
}

func updateLine(content *[]string, replacement string, path ...string) bool {
	idx := findLineIdx(*content, path...)

	if idx == -1 {
		return false
	}

	(*content)[idx] = replacement

	return true
}

func updateVersion(productSpec *string, mp string, version string) bool {
	content := strings.Split(*productSpec, "\n")

	//update toolchains
	updated1 := updateLine(&content, `            "version": "`+version+`"`,
		`"build": \{`, `"toolchains": \{`, `"gradle": \{`, `"plugins": \{`, `"`+mp+`": \{`, `"version": ".*"`)

	//update product dependencies
	updated2 := updateLine(&content, `      "version": "`+version+`"`,
		`"product": \{`, `"`+mp+`": \{`, `"version": ".*"`)

	if updated1 || updated2 {
		*productSpec = strings.Join(content, "\n")
		return true
	} else {
		return false
	}
}
