package gen

import "regexp"

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
