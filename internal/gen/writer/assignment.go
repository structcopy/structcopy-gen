package writer

import (
	"strings"

	"github.com/bookweb/structcopy-gen/internal/gen/writer/model"
)

// AssignmentToString returns the string representation of the assignment.
func AssignmentToString(f *model.Function, a model.Assignment) string {
	var sb strings.Builder
	sb.WriteString(a.String())
	if a.RetError() {
		if f.DstVarStyle == model.DstVarReturn && f.Dst.Pointer {
			sb.WriteString("if err != nil {\nreturn nil, err\n}\n")
		} else {
			sb.WriteString("if err != nil {\nreturn\n}\n")
		}
	}
	return sb.String()
}
