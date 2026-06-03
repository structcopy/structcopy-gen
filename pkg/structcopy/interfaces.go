package structcopy

// Interface represents a complete interface type declaration.
type Interface struct {
	Name         string
	Methods      []Method
	ReceiverType string
	ReceiverName string
}
