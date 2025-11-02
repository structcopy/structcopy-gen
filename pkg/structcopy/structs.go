package structcopy

type Struct struct {
	PackagePath string // where this struct is defined ("" for local)
	PkgName     string
	Name        string // struct name
	Type        string
	Fields      []Field // struct fields
}
