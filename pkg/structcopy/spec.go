package structcopy

// Spec is the root structure to hold all extracted interfaces from a file.
type Spec struct {
	PackageName string
	Imports     []Import
	Interfaces  []Interface
}

// TypeKind enumerates the fundamental types we want to track.
type TypeKind string

const (
	KindBasic     TypeKind = "basic"
	KindStruct    TypeKind = "struct"
	KindInterface TypeKind = "interface"
	KindSlice     TypeKind = "slice"
	KindPointer   TypeKind = "pointer"
	KindMap       TypeKind = "map"
	KindUnknown   TypeKind = "unknown"
)

type TypeInfo struct {
	Name    string
	Kind    TypeKind
	Type    string     // full type
	Fields  []TypeInfo // for struct
	Element *TypeInfo  // for pointer, slice
	Key     *TypeInfo  // for map
	Value   *TypeInfo  // for map
}

// Field represents info about a struct field
type Field struct {
	Name       string // field name
	Kind       string
	Type       string // raw type name (UserID, *User, etc.)
	FullType   string
	IsStruct   bool
	IsPointer  bool   // true if field type is pointer
	IsSlice    bool   // true if field type is slice []User, []*User
	PackageRef string // package import path if external type ("" if local)
}

type MethodParam struct {
	Name                string
	Kind                string
	Type                string
	PointerlessFullType string
	FullType            string
	PackageRef          string // "" if local
	IsStruct            bool
	IsPointer           bool
	IsSlice             bool
	StructDef           *Struct // ParsedStruct if we found its definition
}

type MethodResult struct {
	Name                string
	Kind                string
	Type                string
	PointerlessFullType string
	FullType            string
	PackageRef          string // "" if local
	IsStruct            bool
	IsPointer           bool
	IsSlice             bool
	StructDef           *Struct // ParsedStruct if we found its definition
}

type ParsedMethod struct {
	Name   string
	Params []MethodParam
}
