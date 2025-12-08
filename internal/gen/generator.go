package gen

import (
	"fmt"
	"go/ast"
	"go/token"
	"log/slog"
	"os"
	"strings"
	"unicode"

	"github.com/structcopy/structcopy-gen/pkg/structcopy"
	"golang.org/x/tools/go/packages"
)

// Generator instance.
type Generator struct {
	pkg  *packages.Package
	fset *token.FileSet
	file *ast.File

	spec   *structcopy.Spec
	input  string
	output string
	log    string
	logs   bool
	prints bool
	dryRun bool

	logger *slog.Logger
}

// NewGenerator returns new Generator instance.
func NewGenerator(pkg *packages.Package, fset *token.FileSet, file *ast.File, opts ...GeneratorOption) (*Generator, error) {
	g := &Generator{
		pkg:  pkg,
		fset: fset,
		file: file,
		spec: &structcopy.Spec{},
	}
	g.initDefaults()

	for _, opt := range opts {
		opt(g)
	}

	pkgName := file.Name.Name
	pkgPath := pkg.PkgPath
	g.spec.PackageName = pkgName

	for _, imp := range file.Imports {
		importPath := imp.Path.Value // The quoted import path, e.g., "\"fmt\""

		currentImport := structcopy.Import{
			Path: importPath,
		}

		var alias string
		if imp.Name != nil {
			alias = imp.Name.Name // The import alias if present, e.g., "f" in "import f \"fmt\""
			currentImport.Name = alias
		}

		g.spec.Imports = append(g.spec.Imports, currentImport)
	}

	structs := map[string]*ast.StructType{}
	parsedStructs := map[string]*structcopy.Struct{}

	// ✅ Collect structs from local + imported packages
	for _, file := range pkg.Syntax {
		collectStructs(file, structs, parsedStructs, pkgPath, pkgName)
	}
	for _, imp := range pkg.Imports {
		for _, file := range imp.Syntax {
			collectStructs(file, structs, parsedStructs, pkgPath, pkgName)
		}
	}

	// Traverse the AST
	for _, decl := range file.Decls {
		// Check for a General Declaration (GenDecl)
		// Interfaces, structs, constants, and variables are all GenDecls.
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		// Within a GenDecl, we are interested in "type" declarations.
		// The token.TYPE constant confirms it's a type declaration.
		if genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				// The spec should be a TypeSpec
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// Check if the type is an interface
				interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
				if !ok {
					continue // Skip non-interface types (like structs)
				}

				// We found an interface!
				interfaceName := typeSpec.Name.Name

				if interfaceName != "StructCopyGen" {
					continue // skip interface which name is not equal 'StructCopyGen'
				}
				g.logger.Info(fmt.Sprintf("Valid Interface: %s", interfaceName))

				// Start building a new Interface struct
				currentInterface := structcopy.Interface{
					Name:    typeSpec.Name.Name,
					Methods: []structcopy.Method{},
				}
				// currentInterfaceOptions := &option.Options{}

				if typeSpec.Doc != nil {
					// opts := option.NewOptions()
					// for _, comment := range typeSpec.Doc.List {
					// }
					// currentInterfaceOptions, _ = g.CollectOptions(typeSpec.Doc.List, option.ValidOpsMethod)
				}

				// Iterate over the method list of the interface
				if interfaceType.Methods != nil {
					for _, method := range interfaceType.Methods.List {
						// Method Name
						if len(method.Names) == 0 {
							continue // Skip embedded interfaces
						}
						methodName := method.Names[0].Name

						g.logger.Info(fmt.Sprintf("Valid Method: %s", methodName))

						// Initialize a new Method struct
						currentMethod := structcopy.Method{
							Name:        method.Names[0].Name,
							DstVarStyle: structcopy.DstVarReturn,
						}

						currentMethodOptions := &structcopy.InputOption{}
						var err error
						// Get Documentation Comment
						if method.Doc != nil {
							currentMethodOptions, err = g.CollectOptions(method.Doc.List, ValidOpsMethod)
							if err != nil {
								g.logger.Error("collect options failed", slog.Any("error", err))
							}
							g.logger.Info("Valid annotations")
						}

						currentMethod.SkipFieldsMap = currentMethodOptions.SkipFieldsMap
						currentMethod.MatchFieldsMap = currentMethodOptions.MatchFieldsMap
						currentMethod.MatchMethodsMap = currentMethodOptions.MatchMethodsMap
						currentMethod.ConvertersMap = currentMethodOptions.ConvertersMap
						currentMethod.StructConverterFunc = currentMethodOptions.StructConverterFunc

						// Get the Function Type (*ast.FuncType) of the method
						funcType, ok := method.Type.(*ast.FuncType)
						if !ok {
							g.logger.Info("ERROR: Method type is not *ast.FuncType")
							continue
						}

						// --- Parameters (Inputs) ---
						// paramStrs := []string{}
						if funcType.Params != nil {
							for _, paramField := range funcType.Params.List {
								var params []structcopy.MethodParam
								params = append(params, parseMethodParams(pkgName, paramField, parsedStructs)...)
								currentMethod.Params = params
							}

							if len(currentMethod.Params) > 0 {
								currentMethod.FirstParam = currentMethod.Params[0]
							}
						}

						// --- Results (Outputs) ---
						// resultStrs := []string{}
						if funcType.Results != nil {
							for _, resultField := range funcType.Results.List {
								var results []structcopy.MethodResult
								results = append(results, parseMethodResults(pkgName, resultField, parsedStructs)...)
								currentMethod.Results = results
							}

							if len(currentMethod.Results) > 0 {
								currentMethod.FirstResult = currentMethod.Results[0]
							}
						}

						assignments, err := g.mkMethodAssignments(
							currentMethod.FirstParam,
							currentMethod.FirstResult,
							currentMethod,
						)
						if err != nil {
							g.logger.Error("make assignments failed", slog.Any("error", err))
							return nil, err
						}
						currentMethod.Assignments = assignments

						currentInterface.Methods = append(currentInterface.Methods, currentMethod)
					}
				}

				g.spec.Interfaces = append(g.spec.Interfaces, currentInterface)
			}
		}
	}

	return g, nil
}

// Helper function to extract the type string from an ast.Expr
func extractType(expr ast.Expr) string {
	// We use the basic ast.Inspect for a simple, recursive traversal
	// to find the identifier names that form the type.
	var typeName strings.Builder
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			typeName.WriteString(x.Name)
		case *ast.StarExpr: // For pointer types, e.g., *T
			typeName.WriteString("*")
		case *ast.SelectorExpr: // For qualified types, e.g., fmt.Stringer
			typeName.WriteString(extractType(x.X) + "." + x.Sel.Name)
		}
		return true // Continue inspecting children
	})

	result := typeName.String()
	if result == "" {
		// Handle complex types like "interface{}" or function types (less common in interfaces)
		return fmt.Sprintf("%T", expr)
	}
	return result
}

// typeStringer is a utility to print the raw Go type string from an AST expression.
func typeStringer(expr ast.Expr) string {
	var sb strings.Builder
	fset := token.NewFileSet()
	_ = ast.Fprint(&sb, fset, expr, nil)
	return strings.TrimSpace(sb.String())
}

// extractTypeInfo recursively analyzes an ast.Expr and returns a structured TypeInfo.
func extractTypeInfo(expr ast.Expr, paramName string) structcopy.TypeInfo {
	info := structcopy.TypeInfo{
		Name: paramName,
		Type: typeStringer(expr),
		Kind: structcopy.KindUnknown,
	}

	switch x := expr.(type) {
	case *ast.Ident:
		info.Kind = structcopy.KindBasic
	case *ast.SelectorExpr: // Qualified type (e.g., "io.Reader")
		info.Kind = structcopy.KindBasic
	case *ast.StarExpr: // Pointer type (e.g., "*T")
		info.Kind = structcopy.KindPointer
		// Recursively analyze the pointed-to type
		elementInfo := extractTypeInfo(x.X, "")
		info.Element = &elementInfo
	case *ast.ArrayType: // Slice or Array type (e.g., "[]T")
		info.Kind = structcopy.KindSlice
		// Recursively analyze the element type
		elementInfo := extractTypeInfo(x.Elt, "")
		info.Element = &elementInfo
	case *ast.MapType: // Map type (e.g., "map[K]V")
		info.Kind = structcopy.KindMap
		// Recursively analyze the key and value types
		keyInfo := extractTypeInfo(x.Key, "")
		valueInfo := extractTypeInfo(x.Value, "")
		info.Key = &keyInfo
		info.Value = &valueInfo
	case *ast.StructType: // Inline struct definition
		info.Kind = structcopy.KindStruct
		info.Fields = extractStructFields(x)
	case *ast.InterfaceType: // Interface definition (e.g., "interface{}")
		info.Kind = structcopy.KindInterface
		// We could analyze methods here, but for simplicity, we just mark the kind.
	default:
		// Fallback for types not explicitly covered (FuncType, ChanType, etc.)
		info.Kind = structcopy.KindUnknown
	}

	return info
}

// extractStructFields processes the fields of an *ast.StructType.
func extractStructFields(structType *ast.StructType) []structcopy.TypeInfo {
	var fields []structcopy.TypeInfo
	for _, field := range structType.Fields.List {
		// A struct field can have multiple names (e.g., 'A, B int')
		fieldNames := []string{}
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				fieldNames = append(fieldNames, name.Name)
			}
		} else {
			// Embedded type, use the type name as the field name
			fieldNames = append(fieldNames, typeStringer(field.Type))
		}

		// Recursively call extractTypeInfo for the field's type
		for _, name := range fieldNames {
			fields = append(fields, extractTypeInfo(field.Type, name))
		}
	}
	return fields
}

func (g *Generator) CollectOptions(notations []*ast.Comment, validOps map[string]struct{}) (*structcopy.InputOption, error) {
	inputOption := &structcopy.InputOption{
		StructConverterFunc: "",
		SkipFieldsMap:       map[string]bool{},
		MatchFieldsMap:      map[string]string{},
		MatchMethodsMap:     map[string]string{},
		ConvertersMap:       map[string]string{},
	}

	for _, n := range notations {
		m := reNotation.FindStringSubmatch(n.Text)
		if m == nil || len(m) < 2 {
			return nil, fmt.Errorf("invalid notation format %#v", m)
		}

		var args []string
		if len(m) == 3 {
			args = strings.Fields(m[2])
		}

		if _, ok := validOps[m[1]]; !ok {
			g.logger.Info(fmt.Sprintf(`%v: ":%v" is invalid or unknown notation here`, g.fset.Position(n.Pos()), m[1]))
			continue
		}

		switch m[1] {
		case "structcopygen":
			// do nothing
		case "skip_field":
			if len(args) < 1 {
				return nil, fmt.Errorf("%v: needs <dst> args", g.fset.Position(n.Pos()))
			}
			dst := args[0]

			inputOption.SkipFieldsMap[dst] = true
		case "match_field":
			if len(args) < 2 {
				return nil, fmt.Errorf("%v: needs <dst> <src> args", g.fset.Position(n.Pos()))
			}
			dst := args[0]
			src := args[1]

			inputOption.MatchFieldsMap[dst] = src
		case "match_method":
			if len(args) < 2 {
				return nil, fmt.Errorf("%v: needs <dst> <method> args", g.fset.Position(n.Pos()))
			}
			dst := args[0]
			method := args[1]

			inputOption.MatchMethodsMap[dst] = method
		case "conv":
			if len(args) < 2 {
				return nil, fmt.Errorf("%v: needs <dst> <convert_func> args", g.fset.Position(n.Pos()))
			}
			dst := args[0]
			convertFunc := args[1]

			inputOption.ConvertersMap[dst] = convertFunc
		case "struct_conv":
			if len(args) < 1 {
				return nil, fmt.Errorf("%v: needs <convert_func> args", g.fset.Position(n.Pos()))
			}
			convertFunc := args[0]

			inputOption.StructConverterFunc = convertFunc
		default:
			fmt.Printf("%v: unknown notation %v\n", g.fset.Position(n.Pos()), m[1])
		}
	}

	return inputOption, nil
}

// isValidIdentifier checks if the given string is a valid identifier.
func isValidIdentifier(id string) bool {
	for i, r := range id {
		if !unicode.IsLetter(r) &&
			!(0 < i && unicode.IsDigit(r)) {
			return false
		}
	}
	return id != ""
}

func (g *Generator) initDefaults() {
	g.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func collectStructs(file *ast.File, store map[string]*ast.StructType, parsedStore map[string]*structcopy.Struct, pkgPath, rootPkgName string) {
	packageName := file.Name.Name
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}

		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// structPkgName := ""
			// switch t := ts.Type.(type) {
			// case *ast.Ident:
			// 	structPkgName = packageName
			// case *ast.SelectorExpr:
			// 	pkgIdent, ok := t.X.(*ast.Ident)
			// 	if ok {
			// 		structPkgName = pkgIdent.Name
			// 	}
			// case *ast.StarExpr:
			// }

			if st, ok := ts.Type.(*ast.StructType); ok {
				structPkgName := ""
				if packageName != rootPkgName {
					structPkgName = packageName
				}

				structTypeName := ts.Name.Name
				structMapKey := structTypeName
				if structPkgName != "" {
					structMapKey = fmt.Sprintf("%s.%s", structPkgName, structTypeName)
				}
				store[ts.Name.Name] = st
				parsedStore[structMapKey] = parseStruct(ts, structPkgName, pkgPath)
			}
		}
	}
}

// ---- Helpers ----

// extractTypeName unwraps *T, pkg.T, *pkg.T → "T"
func extractTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident: // User
		return t.Name
	case *ast.StarExpr: // *User, *models.User
		return extractTypeName(t.X)
	case *ast.SelectorExpr: // models.User
		return extractTypeName(t.Sel)
	}
	return ""
}

func stripType(t string) string {
	t = strings.TrimPrefix(t, "*")
	if strings.Contains(t, ".") {
		parts := strings.Split(t, ".")
		return parts[len(parts)-1]
	}
	return t
}

func parseStruct(typeSpec *ast.TypeSpec, pkgName string, rootPkgPath string) *structcopy.Struct {
	st, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	result := structcopy.Struct{
		PackagePath: rootPkgPath,
		PkgName:     pkgName,
		Type:        typeSpec.Name.Name,
		Name:        typeSpec.Name.Name,
	}

	for _, field := range st.Fields.List {
		typeName, pkgRef, isPtr, isSlice := parseFieldType(pkgName, field.Type)

		for _, name := range field.Names { // handle embedded fields too later
			result.Fields = append(result.Fields, structcopy.Field{
				Name:       name.Name,
				Type:       typeName,
				PackageRef: pkgRef,
				IsPointer:  isPtr,
				IsSlice:    isSlice,
			})
		}
	}

	return &result
}

func parseFieldType(pkgName string, expr ast.Expr) (typeName string, pkg string, isPtr, isSlice bool) {
	switch t := expr.(type) {
	case *ast.Ident: // basic type
		return t.Name, "", false, false
	case *ast.SelectorExpr: // imported type
		pkgIdent, ok := t.X.(*ast.Ident)
		if !ok {
			return t.Sel.Name, "", false, false
		}
		return t.Sel.Name, pkgIdent.Name, false, false
	case *ast.StarExpr: // pointer type, includes pointer basic type and pointer imported type
		name, pkg, _, isSliceInner := parseFieldType(pkgName, t.X)
		return name, pkg, true, isSliceInner
	case *ast.ArrayType:
		name, pkg, isPtrInner, _ := parseFieldType("", t.Elt)
		return name, pkg, isPtrInner, true
	// case *ast.MapType: // Map type (e.g., "map[K]V")
	// 	// info.Kind = structcopy.KindMap
	// 	// Recursively analyze the key and value types
	// 	// keyInfo := parseFieldType("", t.Key)
	// 	// valueInfo := parseFieldType("", t.Value)
	// 	// info.Key = &keyInfo
	// 	// info.Value = &valueInfo
	// case *ast.StructType: // Inline struct definition
	// 	// info.Kind = structcopy.KindStruct
	// 	// info.Fields = extractStructFields(t)
	// case *ast.InterfaceType: // Interface definition (e.g., "interface{}")
	// 	// info.Kind = structcopy.KindInterface
	// 	// We could analyze methods here, but for simplicity, we just mark the kind.
	default:
		return fmt.Sprintf("%T", expr), "", false, false
	}
}

func parseMethodParams(pkgName string, field *ast.Field, structs map[string]*structcopy.Struct) []structcopy.MethodParam {
	var results []structcopy.MethodParam

	paramName := "src"
	// (src *entity.User)
	for _, name := range field.Names {
		paramName = name.Name
	}
	typeName, pkgRef, isPointer, isSlice := parseFieldType(pkgName, field.Type)

	key := typeName
	if pkgRef != "" {
		key = pkgRef + "." + typeName
	}

	isStruct := false
	structDef, ok := structs[key]
	if ok {
		isStruct = true
	}

	fullType := key
	if isPointer {
		fullType = fmt.Sprintf("*%s", key)
	}
	param := structcopy.MethodParam{
		Name:                paramName,
		Type:                typeName,
		PointerlessFullType: key,
		FullType:            fullType,
		PackageRef:          pkgRef,
		IsStruct:            isStruct,
		IsPointer:           isPointer,
		IsSlice:             isSlice,
		StructDef:           structDef, // link to collected struct if exists
	}

	results = append(results, param)

	return results
}

func parseMethodResults(pkgName string, field *ast.Field, structs map[string]*structcopy.Struct) []structcopy.MethodResult {
	var results []structcopy.MethodResult

	resultName := "dst"
	// (dst *dto.UserDTO)
	for _, name := range field.Names {
		resultName = name.Name
	}
	typeName, pkgRef, isPtr, isSlice := parseFieldType(pkgName, field.Type)

	key := typeName
	if pkgRef != "" {
		key = pkgRef + "." + typeName
	}

	isStruct := false
	structDef, ok := structs[key]
	if ok {
		isStruct = true
	}

	fullType := key
	if isPtr {
		fullType = fmt.Sprintf("*%s", key)
	}
	param := structcopy.MethodResult{
		Name:                resultName,
		Type:                typeName,
		PointerlessFullType: key,
		FullType:            fullType,
		PackageRef:          pkgRef,
		IsStruct:            isStruct,
		IsPointer:           isPtr,
		IsSlice:             isSlice,
		StructDef:           structDef, // link to collected struct if exists
	}

	results = append(results, param)

	return results
}
