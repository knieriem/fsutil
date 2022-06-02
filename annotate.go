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
func OSName(name string, fsys fs.FS) (string, error) {
	if !fs.ValidPath(name) {
		return "", fs.ErrInvalid
	}
	afs, ok := fsys.(*annotateFS)
	if !ok {
		return "", fs.ErrInvalid
	}
	if name == "." {
		return afs.OSDir, nil
	}
	if pfx := afs.Prefix; pfx != "." {
		name = strings.TrimPrefix(name, pfx)
		if strings.HasPrefix(name, "/") {
			name = name[1:]
		}
	}
	if afs.OSDir != "" {
		name = path.Join(afs.OSDir, name)
	}
	return name, nil
}
