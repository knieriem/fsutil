// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fsutil

import "io/fs"

var _ fs.ReadDirFS = UnionFS{}

// A UnionFS is an FS presenting the union of the file systems in the slice.
// If multiple file systems provide a particular file, Open uses the FS listed earlier in the slice.
// If multiple file systems provide a particular directory, ReadDir presents the
// concatenation of all the directories listed in the slice (with duplicates removed).
type UnionFS []fs.FS

func (fsys UnionFS) Open(name string) (fs.File, error) {
	var errOut error
	for _, sub := range fsys {
		f, err := sub.Open(name)
		if err == nil {
			// Note: Should technically check for directory
			// and return a synthetic directory that merges
			// reads from all the matching directories,
			// but all the directory reads in internal/godoc
			// come from fsys.ReadDir, which does that for us.
			// So we can ignore direct f.ReadDir calls.
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
	for _, sub := range fsys {
		list, err := fs.ReadDir(sub, name)
		if err != nil {
			errOut = err
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
	if len(all) > 0 {
		return all, nil
	}
	return nil, errOut
}
