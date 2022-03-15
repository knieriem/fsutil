package fsutil

import (
	"io/fs"

	golangorg "github.com/knieriem/fsutil/internal/golang-website-golangorg"
)

// A UnionFS is an FS presenting the union of the file systems in a slice.
// If multiple file systems provide a particular file, Open uses the FS listed earlier in the slice.
// If multiple file systems provide a particular directory, ReadDir presents the
// concatenation of all the directories listed in the slice (with duplicates removed).
type UnionFS = golangorg.UnionFS

// Item constitutes an interface to the file system an item
// belongs to. Item is implemented by Files and DirEntries
// provides by a UnionFS.
type Item interface {
	// FS returns the item's file system.
	FS() fs.FS
}
