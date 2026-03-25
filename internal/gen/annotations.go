package gen

import (
	"regexp"
)

var (
	// reNotation is a regular expression that matches a notation.
	reNotation = regexp.MustCompile(`^\s*//\s*:(\S+)\s*(.*)$`)
	// reStructcopygen is a regular expression that matches a notation that
	// indicates the beginning of a structcopy-gen block.
	reStructcopygen = regexp.MustCompile(`^\s*//\s*:structcopygen\b`)
	// reLiteral is a regular expression that matches a notation that
	// indicates the beginning of a literal block.
	reLiteral = regexp.MustCompile(`^\s*\S+\s+(.*)$`)
)

// ValidOpsIntf is a set of valid conversion option keys for interface-level conversion.
var ValidOpsIntf = map[string]struct{}{
	"structcopygen": {},
	"receiver_type": {},
}

// ValidOpsMethod is a set of valid conversion option keys for method-level conversion.
var ValidOpsMethod = map[string]struct{}{
	"skip_field":   {},
	"match_field":  {},
	"match_method": {},
	"conv":         {},
	"struct_conv":  {},
}
