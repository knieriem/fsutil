package fsutil

import (
	"io/fs"
	"path"
	"strings"
	"sync"
)

// annotateFS is a wrapper that allows to store some information
// related to the underlying file system
type annotateFS struct {
	fs.FS

	mu     sync.RWMutex
	values map[interface{}]interface{}
}

func (a *annotateFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(a.FS, name)
}

type ValueKey int

const (
	RootOSDirKey ValueKey = iota
	LabelKey
)

// OSName translates the given file or directory name relative to
// the OS file system tree root, where the file or directory actually resides.
// OSName is able to return the corresponding OS file name only
// if the file system has been bound to a NameSpace
// with option WithNewOSDir previously.
func OSName(name string, fsys fs.FS) (string, error) {
	if !fs.ValidPath(name) {
		return "", fs.ErrInvalid
	}
	afs, ok := fsys.(*annotateFS)
	if !ok {
		return "", fs.ErrInvalid
	}
	afs.mu.RLock()
	defer afs.mu.RUnlock()
	rootOSDir, _ := afs.values[RootOSDirKey].(string)
	if name == "." {
		return rootOSDir, nil
	}
	if fsys, ok := afs.FS.(*prefixFS); ok {
		name = strings.TrimPrefix(name, fsys.pathPrefix)
		if strings.HasPrefix(name, "/") {
			name = name[1:]
		}
	}
	if rootOSDir != "" {
		name = path.Join(rootOSDir, name)
	}
	return name, nil
}

func (a *annotateFS) setValue(key, value interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.values == nil {
		a.values = make(map[interface{}]interface{})
	}
	a.values[key] = value
}

func Value(fsys fs.FS, key interface{}) interface{} {
	afs, ok := fsys.(*annotateFS)
	if !ok {
		return nil
	}

	afs.mu.RLock()
	defer afs.mu.RUnlock()
	return afs.values[key]
}
