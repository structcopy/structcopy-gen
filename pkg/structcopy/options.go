package structcopy

type InterfaceOption struct {
	IsStructCopyGen bool
	ReceiverType    string // n, s, f
	ReceiverName    string // default: myConverter
}

type InputOption struct {
	SkipFieldsMap       map[string]bool
	MatchFieldsMap      map[string]string
	MatchMethodsMap     map[string]string
	ConvertersMap       map[string]string
	StructConverterFunc string
}
