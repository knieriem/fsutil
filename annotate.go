package fsutil

import (
	"io/fs"
	"path"
	"strings"
)

// annotateFS is a wrapper that allows to store some information
// related to the underlying file system
type annotateFS struct {
	fs.FS
	Label  string
	Prefix string
	OSDir  string
}

func (a *annotateFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(a.FS, name)
}

// OSName translates the given file or directory name relative to
// the OS file system tree root, where the file or directory actually resides.
// OSName is able to return the corresponding OS file name only
// if the file system has been bound to a NameSpace
// with option WithNewOSDir previously.
func OSName(name string, fsys fs.FS) (osName string, ok bool) {
	afs, ok := fsys.(*annotateFS)
	if ok {
		name = strings.TrimPrefix(name, afs.Prefix)
		if afs.OSDir != "" {
			return path.Join(afs.OSDir, name), true
		}
	}
	return "", false
}
