This repository contains utilities around Go's `io/fs` package:

-	UnionFS, which is a copy of the [`unionFS` type implemented
	in golang.org/x/website/cmd/golangorg/server.go][unionFS],
	slightly extended to allow access to the underlying
	file system of Files and DirEntrys

-	NameSpace, a wrapper around UnionFS that provides
	method Bind, similar to [Plan 9's bind] command.
	It has been created to be able to use a functionality like Godoc's [vfs.NameSpace],
	but based on `io/fs`.

-	PrefixFS creates a filesystem by adding a prefix in front of another,
	already existing FS; it is used by _Namespace._

[Plan 9's bind]: https://9p.io/magic/man2html/1/bind
[unionFS]: https://cs.opensource.google/go/x/website/+/master:cmd/golangorg/server.go;l=591-649
[vfs.NameSpace]: https://pkg.go.dev/golang.org/x/tools/godoc/vfs#NameSpace
