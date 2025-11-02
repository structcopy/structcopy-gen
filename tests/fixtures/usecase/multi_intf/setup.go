package multi_intf

type DomainModel struct {
	ID string
}

type TransportModel struct {
	ID string
}

type StorageModel struct {
	ID string
}

//go:generate structcopy-gen
type StructCopyGen interface {
	// :recv d
	ToTransport(*DomainModel) *TransportModel
	// :recv d
	ToStorage(*DomainModel) *StorageModel
}

// :structcopygen
type StorageConverter interface {
	// :recv s
	ToTransport(*StorageModel) *TransportModel
	// :recv s
	ToDomain(*StorageModel) *DomainModel
}
