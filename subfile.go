package fsutil

import (
	"fmt"
	"io/fs"
	"path"
)

// SubFile returns an [fs.FS] rooted at the directory containing
// the named file, with that file as the only entry.
// SubFile does not check whether the file exists.
func SubFile(fsys fs.FS, name string) (fs.FS, error) {
	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("fsutil.SubFile: invalid path %q", name)
	}
	sub := new(subFileFS)
	sub.fsys = fsys
	sub.origName = name
	sub.name = path.Base(name)
	return sub, nil
}

type subFileFS struct {
	fsys     fs.FS
	name     string
	origName string
}

func (sub *subFileFS) Open(name string) (fs.File, error) {
	if name != sub.name {
		return nil, fs.ErrNotExist
	}
	return sub.fsys.Open(sub.origName)
}

func (sub *subFileFS) Stat(name string) (fs.FileInfo, error) {
	if name == "." {
		// Calling Stat on the underlying file system
		// as a hack to make fs.WalkDir succeed.
		return fs.Stat(sub.fsys, path.Dir(sub.origName))
	}
	if name != sub.name {
		return nil, fs.ErrNotExist
	}
	return (&subFileDir{sub}).Info()
}

func (sub *subFileFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var dir [1]fs.DirEntry

	dir[0] = &subFileDir{sub}
	return dir[:], nil
}

type subFileDir struct {
	*subFileFS
}

func (d *subFileDir) Name() string {
	return d.name
}

func (d *subFileDir) IsDir() bool {
	return false
}

func (d *subFileDir) Type() fs.FileMode {
	return 0
}

func (d *subFileDir) Info() (fs.FileInfo, error) {
	return fs.Stat(d.fsys, d.origName)
}
