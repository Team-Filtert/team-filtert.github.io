package main

import (
	"path/filepath"
	"strings"
)

// stringList is a repeatable flag value: --changed a --changed b.
type stringList []string

func (s *stringList) String() string {
	return strings.Join(*s, ",")
}

func (s *stringList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var incrementalPrefixes = []string{"content/", "public/", "style/", "scripts/"}

func isIncrementalPath(path string) bool {
	for _, prefix := range incrementalPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// dataOnly reports whether every changed path can be handled by an incremental
// build. Go and templ sources need a full build.
func dataOnly(changed []string) bool {
	for _, path := range changed {
		if !isIncrementalPath(filepath.ToSlash(path)) {
			return false
		}
	}
	return true
}
