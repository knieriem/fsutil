// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fsutil

import (
	"io"
	"io/fs"
)

var _ fs.ReadDirFS = UnionFS{}

// A UnionFS is an FS presenting the union of the file systems in the slice.
// If multiple file systems provide a particular file, Open uses the FS listed earlier in the slice.
// If multiple file systems provide a particular directory, ReadDir presents the
// concatenation of all the directories listed in the slice (with duplicates removed).
type UnionFS []fs.FS

func (fsys UnionFS) Open(name string) (fs.File, error) {
	var errOut error
	if len(fsys) == 0 {
		return nil, fs.ErrNotExist
	}
	for _, sub := range fsys {
		f, err := sub.Open(name)
		if err == nil {
			fi, err := f.Stat()
			if err != nil {
				return nil, err
			}
			if fi.IsDir() {
				return &dir{file: file{File: f, fsys: sub}, name: name, union: fsys}, nil
			}
			return &file{File: f, fsys: sub}, nil
		}
		if errOut == nil {
			errOut = err
		}
	}
	return nil, errOut
}

type file struct {
	fs.File
	fsys fs.FS
}

func (f file) FS() fs.FS {
	return f.fsys
}

type dir struct {
	file
	name  string
	union fs.ReadDirFS
	list  []fs.DirEntry
}

func (d *dir) ReadDir(n int) ([]fs.DirEntry, error) {
	if d.list == nil {
		list, err := d.union.ReadDir(d.name)
		if err != nil {
			return nil, err
		}
		if n <= 0 {
			return list, nil
		}
		d.list = list
	}
	if len(d.list) == 0 {
		return nil, io.EOF
	}
	if n > len(d.list) {
		n = len(d.list)
	}
	list := d.list[:n]
	d.list = d.list[n:]
	return list, nil
}

type dirEntry struct {
	fs.DirEntry
	fsys fs.FS
}

func (d dirEntry) FS() fs.FS {
	return d.fsys
}

func (fsys UnionFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var all []fs.DirEntry
	var seen map[string]bool // seen[name] is true if name is listed in all; lazily initialized
	var errOut error
	anyOK := false
	for _, sub := range fsys {
		list, err := fs.ReadDir(sub, name)
		if err != nil {
			errOut = err
		} else {
			anyOK = true
		}
		if len(list) == 0 {
			continue
		}
		if len(all) == 0 {
			all = make([]fs.DirEntry, len(list))
			for i, d := range list {
				all[i] = &dirEntry{DirEntry: d, fsys: sub}
			}
		} else {
			if seen == nil {
				// Initialize seen only after we get two different directory listings.
				seen = make(map[string]bool)
				for _, d := range all {
					seen[d.Name()] = true
				}
			}
			for _, d := range list {
				name := d.Name()
				if !seen[name] {
					seen[name] = true
					all = append(all, &dirEntry{DirEntry: d, fsys: sub})
				}
			}
		}
	}
	if len(all) > 0 || anyOK {
		return all, nil
	}
	return nil, errOut
}
