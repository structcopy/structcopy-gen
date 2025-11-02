//go:build structcopygen

package multi_intf

// :structcopygen
// TransportConverter should not be handled by the generator since
// the generator treats only setup.go, the "input" file as the StructCopyGen definition file.
type TransportConverter interface {
	// :recv t
	ToDomain(*TransportModel) *DomainModel
	// :recv t
	ToStorage(*TransportModel) *StorageModel
}

// FooBar can be reached from functions in setup.go during the parse phase.
// However, DON'T CALL THIS.
// While no code here won't be included to the generated code, it causes compile errors.
func FooBar() {}
