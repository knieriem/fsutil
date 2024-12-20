package fsutil

import (
	"io"
	"io/fs"
	"slices"
	"strings"
	"time"

	"maps"
)

// StringMap returns an [fs.FS] consisting of the values of the specified map,
// with keys used as filenames. For now, only a plain structure is allowed,
// i.e. filenames must not contain slashes.
func StringMap(m map[string]string) fs.FS {
	filenames := slices.Sorted(maps.Keys(m))
	dir := make([]fs.DirEntry, 0, len(filenames))
	for _, name := range filenames {
		if name == "." {
			continue
		}
		if strings.ContainsRune(name, '/') {
			continue
		}
		d := new(constDirEnt)
		d.name = name
		d.size = len(m[name])
		dir = append(dir, d)
	}
	fsys := &stringMapFS{
		dir: dir,
		m:   m,
	}
	return fsys
}

type stringMapFS struct {
	dir []fs.DirEntry
	m   map[string]string
}

type constDirEnt struct {
	name string
	size int
}

func (f *constDirEnt) Name() string {
	return f.name
}

func (f *constDirEnt) IsDir() bool {
	return false
}

func (f *constDirEnt) Type() fs.FileMode {
	return 0
}

func (f *constDirEnt) Mode() fs.FileMode {
	return 0644
}

func (f *constDirEnt) ModTime() time.Time {
	var zero time.Time
	return zero
}

func (f *constDirEnt) Sys() any {
	return nil
}

func (f *constDirEnt) Info() (fs.FileInfo, error) {
	return f, nil
}

func (f *constDirEnt) Size() int64 {
	return int64(f.size)
}

type constFile struct {
	constDirEnt
	io.Reader
}

func (f *constFile) Stat() (fs.FileInfo, error) {
	return &f.constDirEnt, nil
}

func (f *constFile) Close() error {
	return nil
}

func (fsys *stringMapFS) Open(name string) (fs.File, error) {
	if name == "." {
		return &stringMapDir{dir: fsys.dir}, nil
	}
	v, ok := fsys.m[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	f := new(constFile)
	f.name = name
	f.size = len(v)
	f.Reader = strings.NewReader(v)
	return f, nil
}

func (fsys *stringMapFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if name != "." {
		return nil, fs.ErrNotExist
	}
	return fsys.dir, nil
}

type stringMapDir struct {
	dir []fs.DirEntry
}

func (d *stringMapDir) Read([]byte) (int, error) {
	return 0, fs.ErrPermission
}

func (d *stringMapDir) Stat() (fs.FileInfo, error) {
	return &dummyRoot{}, nil
}

func (d *stringMapDir) Close() error { return nil }

func (d *stringMapDir) ReadDir(n int) ([]fs.DirEntry, error) {
	dir := d.dir
	if dir == nil {
		return nil, io.EOF
	}
	if n <= 0 {
		d.dir = nil
		return dir, nil
	}
	if n > len(dir) {
		n = len(dir)
	}
	dir = dir[:n]
	d.dir = dir[n:]
	return dir, nil
}

type dummyRoot struct {
	constDirEnt
}

func (r *dummyRoot) Name() string { return "." }
func (r *dummyRoot) Size() int64  { return 0 }
func (r *dummyRoot) IsDir() bool  { return true }
