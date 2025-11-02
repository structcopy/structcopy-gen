//go:build structcopygen

package slice

type SrcType struct {
	IntSlice    []int
	DataSlice   []Data
	StatusSlice []int
}

type DstType struct {
	IntSlice    []int
	DataSlice   []Data
	StatusSlice []Status
}

type Data struct{}

type Status int

//go:generate structcopy-gen
type StructCopyGen interface {
	// :typecast
	Copy(*SrcType) *DstType
}
