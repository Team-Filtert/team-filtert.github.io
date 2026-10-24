package main

import (
	"io/fs"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"
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

type fileState struct {
	modTime time.Time
	size    int64
}

// changeTracker remembers the source files from the previous build so the dev
// server can pass an incremental changed set instead of rebuilding everything.
type changeTracker struct {
	previous map[string]fileState
}

func newChangeTracker() *changeTracker {
	return &changeTracker{}
}

// diff returns the files that changed since the last build. initial is true
// for the first build, removed is true when a source disappeared, and changed
// holds absolute slash-separated paths.
func (t *changeTracker) diff() (changed []string, removed, initial bool, err error) {
	current, err := scanFiles()
	if err != nil {
		return nil, false, false, err
	}

	if t.previous == nil {
		t.previous = current
		return nil, false, true, nil
	}

	for path, state := range current {
		previous, ok := t.previous[path]
		if !ok || previous.modTime != state.modTime || previous.size != state.size {
			changed = append(changed, path)
		}
	}
	for path := range t.previous {
		if _, ok := current[path]; !ok {
			removed = true
			break
		}
	}

	t.previous = current
	sort.Strings(changed)
	return changed, removed, false, nil
}

// dataOnly reports whether every changed file can be handled by an incremental
// build. Go and templ sources need a full build because they change the
// compiled pages themselves.
func (t *changeTracker) dataOnly(changed []string) bool {
	root, err := filepath.Abs(".")
	if err != nil {
		return false
	}

	for _, path := range changed {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return false
		}
		rel = filepath.ToSlash(rel)
		if !isIncrementalPath(rel) {
			return false
		}
	}
	return true
}

var incrementalPrefixes = []string{"content/", "public/", "style/", "scripts/"}

func isIncrementalPath(path string) bool {
	return slices.ContainsFunc(incrementalPrefixes, func(prefix string) bool {
		return strings.HasPrefix(path, prefix)
	})
}

func scanFiles() (map[string]fileState, error) {
	root, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}

	states := make(map[string]fileState)
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root && ignoredDir(entry.Name()) {
				return fs.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(entry.Name(), "_templ.go") {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		states[filepath.ToSlash(path)] = fileState{
			modTime: info.ModTime(),
			size:    info.Size(),
		}
		return nil
	})
	return states, err
}

func ignoredDir(name string) bool {
	switch name {
	case ".git", "node_modules", "vendor", ".cache", "build", ".gomodcache", ".gocache":
		return true
	}
	return false
}
