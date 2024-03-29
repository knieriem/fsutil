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
// otherwise after possibly existing files within `old`.
func (nsp *NameSpace) Bind(old string, newfs fs.FS, options ...BindOption) error {
	if !fs.ValidPath(old) {
		return fs.ErrInvalid
	}

	var a bindAction
	for _, o := range options {
		o(&a)
	}
	if old != "." {
		newfs = PrefixFS(old, newfs)
	}
	a.fs.FS = newfs
	if a.newOSDir != "" {
		a.fs.setValue(RootOSDirKey, a.newOSDir)
	}

	if !a.before || len(nsp.UnionFS) == 0 {
		nsp.append(&a.fs)
		return nil
	}
	nsp.UnionFS = append(nsp.UnionFS[:1], nsp.UnionFS...)
	nsp.UnionFS[0] = &a.fs
	return nil
}

type BindOption func(*bindAction)

type bindAction struct {
	before   bool
	newOSDir string
	fs       annotateFS
}

func BindBefore() BindOption {
	return func(a *bindAction) {
		a.before = true
	}
}

func WithValue(key, value interface{}) BindOption {
	return func(a *bindAction) {
		a.fs.setValue(key, value)
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
