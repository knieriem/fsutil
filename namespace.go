// Package fsutil provides utilities around package io/fs,
// mainly a namespace feature built upen a union file system.
package fsutil

import "io/fs"

// Namespace is a wrapper around UnionFS that provides
// a functionality basically similar to Plan 9's bind command
type NameSpace struct {
	UnionFS
}

// Bind adds file system newfs at mount point `old`.
// If option BindBefore is provided, newfs' contents appear first in the union,
// otherweise after possibly existing files within `old`.
func (nsp *NameSpace) Bind(old string, newfs fs.FS, options ...BindOption) {
	var a bindAction
	for _, o := range options {
		o(&a)
	}
	if old != "" {
		newfs = PrefixFS(old, newfs)
	}
	var afs annotateFS
	afs.Prefix = old
	afs.FS = newfs
	afs.OSDir = a.newOSDir

	if !a.before || len(nsp.UnionFS) == 0 {
		nsp.append(&afs)
		return
	}
	nsp.UnionFS = append(nsp.UnionFS[:1], nsp.UnionFS...)
	nsp.UnionFS[0] = &afs
}

type BindOption func(*bindAction)

type bindAction struct {
	before   bool
	newOSDir string
}

func BindBefore() BindOption {
	return func(a *bindAction) {
		a.before = true
	}
}

func WithNewOSDir(dirname string) BindOption {
	return func(a *bindAction) {
		a.newOSDir = dirname
	}
}

func (nsp *NameSpace) append(fsys fs.FS) {
	nsp.UnionFS = append(nsp.UnionFS, fsys)
}
