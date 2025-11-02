package gen

import (
	"errors"
	"fmt"

	"github.com/bookweb/structcopy-gen/pkg/structcopy"
	"github.com/samber/lo"
)

func (g *Generator) mkMethodAssignments(src structcopy.MethodParam, dst structcopy.MethodResult, method structcopy.Method) ([]structcopy.Assignment, error) {
	assignments := make([]structcopy.Assignment, 0)
	if src.IsSlice && dst.IsSlice && src.IsStruct && dst.IsStruct {
		structAssignments, err := g.mkSliceOfStructToSliceOfStructAssignments(src, dst, method)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, structAssignments...)
	} else if src.IsStruct && dst.IsStruct && src.StructDef != nil && dst.StructDef != nil {
		structAssignments, err := g.mkStructToStructAssignments(src, dst, method)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, structAssignments...)
	} else {
		// skip
	}

	return assignments, nil
}

func (g *Generator) mkStructToStructAssignments(
	src structcopy.MethodParam,
	dst structcopy.MethodResult,
	method structcopy.Method,
) ([]structcopy.Assignment, error) {
	assignments := make([]structcopy.Assignment, 0)

	for _, field := range dst.StructDef.Fields {
		assignment, err := g.mkFieldAssignment(field, src, dst, method)
		if err != nil {
			fmt.Println(err)
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

func (g *Generator) mkFieldAssignment(
	field structcopy.Field,
	src structcopy.MethodParam,
	dst structcopy.MethodResult,
	method structcopy.Method,
) (structcopy.Assignment, error) {
	matchFieldsMap := method.MatchFieldsMap
	convertersMap := method.ConvertersMap
	matchMethodsMap := method.MatchMethodsMap

	srcFieldName := field.Name
	matchSrcFieldName, ok := matchFieldsMap[field.Name]
	if ok {
		srcFieldName = matchSrcFieldName
	}

	srcConverterFunc := ""
	converter, ok := convertersMap[field.Name]
	if ok {
		srcConverterFunc = converter
	}

	srcMatchMethod := ""
	matchMethod, ok := matchMethodsMap[field.Name]
	if ok {
		srcMatchMethod = matchMethod
	}

	matchSrcField := lo.ContainsBy(src.StructDef.Fields, func(fi structcopy.Field) bool {
		return fi.Name == field.Name
	})

	lhs := fmt.Sprintf("%s.%s", dst.Name, field.Name)
	rhs := fmt.Sprintf("%s.%s", src.Name, field.Name)
	if srcFieldName != "" {
		rhs = fmt.Sprintf("%s.%s", src.Name, srcFieldName)
	}

	if srcMatchMethod != "" {
		return &structcopy.MatchMethodField{
			LHS:         lhs,
			RContainer:  src.Name,
			MatchMethod: srcMatchMethod,
		}, nil
	} else if srcMatchMethod == "" && matchSrcFieldName == "" && !matchSrcField {
		return &structcopy.NoMatchField{
			LHS: lhs,
		}, nil
	} else if srcConverterFunc != "" {
		return &structcopy.ConvertField{
			LHS:     lhs,
			RHS:     rhs,
			Convert: srcConverterFunc,
		}, nil
	} else {
		return &structcopy.SimpleField{
			LHS: lhs,
			RHS: rhs,
		}, nil
	}
}

func (g *Generator) mkSliceOfStructToSliceOfStructAssignments(
	src structcopy.MethodParam,
	dst structcopy.MethodResult,
	method structcopy.Method,
) ([]structcopy.Assignment, error) {
	structConverterFunc := method.StructConverterFunc
	if structConverterFunc == "" {
		return nil, errors.New("struct_conv func is required")
	}

	assignments := make([]structcopy.Assignment, 0)

	assignment := &structcopy.SliceStructConvertLoopAssignment{
		LHS:           dst.Name,
		RHS:           src.Name,
		Typ:           dst.FullType,
		StructConvert: structConverterFunc,
	}
	assignments = append(assignments, assignment)

	return assignments, nil
}
