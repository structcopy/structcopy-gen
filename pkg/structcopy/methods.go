package structcopy

import (
	"regexp"
	"strings"
)

// Method represents a single function signature within an interface.
type Method struct {
	Name        string
	Params      []MethodParam
	Results     []MethodResult
	FirstParam  MethodParam
	FirstResult MethodResult

	AdditionalArgs []Variable

	Receiver            string
	SkipFieldsMap       map[string]bool
	MatchFieldsMap      map[string]string
	MatchMethodsMap     map[string]string
	ConvertersMap       map[string]string
	StructConverterFunc string
	Assignments         []Assignment
	PreProcess          *Manipulator
	PostProcess         *Manipulator

	Docs     []string
	Comments []string

	DstVarStyle DstVarStyle
	RetError    bool
}

// DstVarStyle represents the style of destination variable in a function signature.
type DstVarStyle string

// String returns the string representation of the destination variable style.
func (s DstVarStyle) String() string {
	return string(s)
}

const (
	// DstVarReturn indicates that the destination variable is a return
	// value in a function signature.
	DstVarReturn = DstVarStyle("return")
	// DstVarArg indicates that the destination variable is an argument
	// in a function signature.
	DstVarArg = DstVarStyle("arg")
)

// DstVarStyleValues is a slice of all possible destination variable styles.
var DstVarStyleValues = []DstVarStyle{DstVarReturn, DstVarArg}

// NewDstVarStyleFromValue creates a new DstVarStyle instance from the
// given value string.
func NewDstVarStyleFromValue(v string) (DstVarStyle, bool) {
	for _, style := range DstVarStyleValues {
		if style.String() == v {
			return style, true
		}
	}
	return "", false
}

// MatchRule represents the field matching rule.
type MatchRule string

// String returns the string representation of the match rule.
func (s MatchRule) String() string {
	return string(s)
}

const (
	// MatchRuleName indicates that the field name is used as the matching criteria.
	MatchRuleName = MatchRule("name")
	// MatchRuleTag indicates that the field tag is used as the matching criteria.
	MatchRuleTag = MatchRule("tag")
	// MatchRuleNone indicates that there is no matching criteria for the field.
	MatchRuleNone = MatchRule("none")
)

// MatchRuleValues is a slice of all possible field matching rules.
var MatchRuleValues = []MatchRule{MatchRuleName, MatchRuleTag, MatchRuleNone}

// NewMatchRuleFromValue creates a new MatchRule instance from the given value string.
func NewMatchRuleFromValue(v string) (MatchRule, bool) {
	for _, rule := range MatchRuleValues {
		if rule.String() == v {
			return rule, true
		}
	}
	return "", false
}

func (f Method) String() string {
	var sb strings.Builder

	// doc comment
	for i := range f.Comments {
		sb.WriteString(f.Comments[i])
		sb.WriteString("\n")
	}

	// "func"
	sb.WriteString("func ")

	// "func (r *SrcModel) Name("
	sb.WriteString(f.Name)
	sb.WriteString("(")

	if f.Receiver == "" {
		// // "func Name(dst *DstModel, src *SrcModel"
		// sb.WriteString(f.Src.Name)
		// sb.WriteString(" ")
		// sb.WriteString(f.Src.FullType())

		// "func Name(dst *DstModel, src *SrcModel"
		sb.WriteString(f.FirstParam.Name)
		sb.WriteString(" ")
		sb.WriteString(f.FirstParam.FullType)
	}

	for _, args := range f.AdditionalArgs {
		fullType := args.FullType()
		if strings.Contains(args.Type, "/") {
			re := regexp.MustCompile(`^([^a-zA-Z0-9]*)([a-zA-Z0-9].*/)(.+)$`)
			fullType = re.ReplaceAllString(fullType, "$1$3")
		}
		sb.WriteString(", ")
		sb.WriteString(args.Name)
		sb.WriteString(" ")
		sb.WriteString(fullType)
	}

	// "func Name(dst *DstModel, src *SrcModel)"
	sb.WriteString(") ")

	if f.DstVarStyle == DstVarReturn {
		// "func Name(src *SrcModel) (dst *DstModel"
		sb.WriteString("(")
		sb.WriteString(f.FirstResult.Name)
		sb.WriteString(" ")
		sb.WriteString(f.FirstResult.FullType)
		if f.RetError {
			// "func Name(src *SrcModel) (dst *DstModel, err error"
			sb.WriteString(", err error")
		}

		// "func Name(src *SrcModel) (dst *DstModel) {"
		sb.WriteString(") {\n")
		if f.FirstResult.IsPointer {
			// "dst = &DstModel{}"
			sb.WriteString(f.FirstResult.Name)
			sb.WriteString(" = ")
			if f.FirstResult.IsPointer {
				sb.WriteString("&")
			}
			sb.WriteString(f.FirstResult.PointerlessFullType)
			sb.WriteString("{}\n")
		}
	} else {
		if f.RetError {
			// "func Name(dst *DstModel, src *SrcModel) (err error) {"
			sb.WriteString("(err error) {\n")
		} else {
			// "func Name(dst *DstModel, src *SrcModel) {"
			sb.WriteString("{\n")
		}
	}

	for i := range f.Assignments {
		sb.WriteString(f.AssignmentToString(f.Assignments[i]))
	}
	if f.RetError || f.DstVarStyle == DstVarReturn {
		sb.WriteString("\nreturn\n")
	}
	sb.WriteString("}\n\n")
	return sb.String()
}

func (f Method) FormatSliceOfStruct() string {
	var sb strings.Builder

	// doc comment
	for i := range f.Comments {
		sb.WriteString(f.Comments[i])
		sb.WriteString("\n")
	}

	// "func"
	sb.WriteString("func ")

	// "func (r *SrcModel) Name("
	sb.WriteString(f.Name)
	sb.WriteString("(")

	if f.Receiver == "" {
		// // "func Name(dst *DstModel, src *SrcModel"
		// sb.WriteString(f.Src.Name)
		// sb.WriteString(" ")
		// sb.WriteString(f.Src.FullType())

		// "func Name(dst *DstModel, src *SrcModel"
		sb.WriteString(f.FirstParam.Name)
		sb.WriteString(" []")
		sb.WriteString(f.FirstParam.FullType)
	}

	for _, args := range f.AdditionalArgs {
		fullType := args.FullType()
		if strings.Contains(args.Type, "/") {
			re := regexp.MustCompile(`^([^a-zA-Z0-9]*)([a-zA-Z0-9].*/)(.+)$`)
			fullType = re.ReplaceAllString(fullType, "$1$3")
		}
		sb.WriteString(", ")
		sb.WriteString(args.Name)
		sb.WriteString(" ")
		sb.WriteString(fullType)
	}

	// "func Name(dst *DstModel, src *SrcModel)"
	sb.WriteString(") ")

	if f.DstVarStyle == DstVarReturn {
		// "func Name(src *SrcModel) (dst *DstModel"
		sb.WriteString("(")
		sb.WriteString(f.FirstResult.Name)
		sb.WriteString(" []")
		sb.WriteString(f.FirstResult.FullType)
		if f.RetError {
			// "func Name(src *SrcModel) (dst *DstModel, err error"
			sb.WriteString(", err error")
		}

		// "func Name(src *SrcModel) (dst *DstModel) {"
		sb.WriteString(") {\n")
		// if f.NewFirstResult.IsPointer {
		// 	// "dst = &DstModel{}"
		// 	sb.WriteString(f.NewFirstResult.Name)
		// 	sb.WriteString(" = ")
		// 	if f.NewFirstResult.IsPointer {
		// 		sb.WriteString("&")
		// 	}
		// 	sb.WriteString(f.NewFirstResult.PointerlessFullType)
		// 	sb.WriteString("{}\n")
		// }
	} else {
		if f.RetError {
			// "func Name(dst *DstModel, src *SrcModel) (err error) {"
			sb.WriteString("(err error) {\n")
		} else {
			// "func Name(dst *DstModel, src *SrcModel) {"
			sb.WriteString("{\n")
		}
	}

	for i := range f.Assignments {
		sb.WriteString(f.AssignmentToString(f.Assignments[i]))
	}

	if f.RetError || f.DstVarStyle == DstVarReturn {
		sb.WriteString("\nreturn\n")
	}
	sb.WriteString("}\n\n")
	return sb.String()
}

// AssignmentToString returns the string representation of the assignment.
func (f Method) AssignmentToString(a Assignment) string {
	var sb strings.Builder
	sb.WriteString(a.String())
	if a.RetError() {
		sb.WriteString("if err != nil {\nreturn\n}\n")
	}
	return sb.String()
}

// ManipulatorToString returns a string representation of the given Manipulator.
// It generates a function call that performs the manipulation and returns the result as a string.
// Parameters:
// - m: the Manipulator to be converted into a string representation.
// - src: the source Var that corresponds to the Manipulator's first argument.
// - dst: the destination Var that corresponds to the Manipulator's second argument.
// Returns:
// - a string that represents the function call to the Manipulator.
func (f Method) ManipulatorToString(m *Manipulator, src, dst Variable, args []Variable) string {
	var sb strings.Builder
	if m.RetError {
		sb.WriteString("err = ")
	}
	if m.Pkg != "" {
		sb.WriteString(m.Pkg)
		sb.WriteString(".")
	}
	sb.WriteString(m.Name)
	sb.WriteString("(")

	if dst.Pointer != m.IsDstPtr {
		if dst.Pointer {
			sb.WriteString("*")
		} else {
			sb.WriteString("&")
		}
	}
	sb.WriteString(dst.Name)
	sb.WriteString(", ")

	if src.Pointer != m.IsSrcPtr {
		if src.Pointer {
			sb.WriteString("*")
		} else {
			sb.WriteString("&")
		}
	}
	sb.WriteString(src.Name)

	if m.HasAdditionalArgs {
		for _, arg := range args {
			sb.WriteString(", ")
			sb.WriteString(arg.Name)
		}
	}
	sb.WriteString(")\n")

	if m.RetError {
		sb.WriteString("if err != nil {\nreturn\n}\n")
	}

	return sb.String()
}
