package structcopy

type InputOption struct {
	MatchFieldsMap      map[string]string
	MatchMethodsMap     map[string]string
	ConvertersMap       map[string]string
	StructConverterFunc string
}
