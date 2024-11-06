package utils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// FunctionInfo holds information about a function and its associated comments.
type FunctionInfo struct {
	FileName string
	Function string
	Comments string
	Params   string // Add a field for the parameters
	Returns  string // Add a field for the return types
}

// PackageFunctions lists and sorts the functions in a package, excluding utils.go and utils_test.go.
func PackageFunctions(dir string) (map[string][]FunctionInfo, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	funcMap := make(map[string][]FunctionInfo)

	for _, pkg := range pkgs {
		for filePath, f := range pkg.Files {
			// Skip utils.go and utils_test.go
			fileName := filepath.Base(filePath)
			if fileName == "utils.go" || fileName == "utils_test.go" {
				continue
			}

			// Use the file name without the extension for grouping
			groupName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

			for _, decl := range f.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.IsExported() {
					// Get the comments associated with the function
					var comments string
					if fn.Doc != nil {
						comments = Commentify(fn.Doc.Text())
					}

					// Extract the function's parameters and return types
					params := extractFuncParams(fn.Type.Params)
					returns := extractFuncResults(fn.Type.Results)

					funcInfo := FunctionInfo{
						FileName: groupName,
						Function: fn.Name.Name,
						Comments: comments,
						Params:   params,
						Returns:  returns,
					}
					funcMap[groupName] = append(funcMap[groupName], funcInfo)
				}
			}
		}
	}

	// Sort the functions within each group alphabetically
	for _, functions := range funcMap {
		sort.Slice(functions, func(i, j int) bool {
			return functions[i].Function < functions[j].Function
		})
	}

	return funcMap, nil
}

// extractFuncParams takes an *ast.FieldList (parameters) and returns a string representation.
func extractFuncParams(fl *ast.FieldList) string {
	if fl == nil {
		return ""
	}
	var params []string
	for _, field := range fl.List {
		typeString := exprToString(field.Type)
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				params = append(params, fmt.Sprintf("%s %s", name, typeString))
			}
		} else {
			params = append(params, typeString)
		}
	}
	return strings.Join(params, ", ")
}

// extractFuncResults takes an *ast.FieldList (results) and returns a string representation.
func extractFuncResults(fl *ast.FieldList) string {
	if fl == nil {
		return ""
	}
	var results []string
	for _, field := range fl.List {
		typeString := exprToString(field.Type)
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				results = append(results, fmt.Sprintf("%s %s", name, typeString))
			}
		} else {
			results = append(results, typeString)
		}
	}
	if len(results) == 1 {
		return results[0] // Single unnamed return value
	}
	return "(" + strings.Join(results, ", ") + ")"
}

// exprToString takes an ast.Expr (expression) and returns its string representation.
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.InterfaceType:
		if t.Methods != nil && t.Methods.List != nil && len(t.Methods.List) > 0 {
			return "interface{...}"
		}
		return "interface{}"
	default:
		return fmt.Sprintf("%#v", expr)
	}
}

// PrintSortedFunctions prints the sorted functions as specified.
func ListPackageFunctions(dir string, packages bool, packageName string) {
	groupedFunctions, err := PackageFunctions(dir)
	if err != nil {
		panic(err)
	}

	var groupNames []string
	for groupName := range groupedFunctions {
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)

	cleanDir := path.Clean(dir)
	lastPart := path.Base(cleanDir)

	if packages {
		header := fmt.Sprintf("\nAvailable Packages:%s", "")
		TerminalColor(header, Info)
		for _, groupName := range groupNames {
			if !strings.Contains(groupName, "_test") {
				TerminalColor("  "+groupName, Notice)
			}
		}
	} else {
		if packageName != "" {
			_, ok := groupedFunctions[packageName]
			if ok {
				header := fmt.Sprintf("\nFunctions in the %s package:", fmt.Sprintf("'%s'", lastPart))
				TerminalColor(header, Info)
				functions := groupedFunctions[packageName]
				fmt.Println()
				TerminalColor("[ "+packageName+" ]", Notice)
				for _, funcInfo := range functions {
					TerminalColor("  "+strings.TrimSpace(funcInfo.Comments), Info)
					signature := fmt.Sprintf("  %s(%s) %s", funcInfo.Function, funcInfo.Params, funcInfo.Returns)
					TerminalColor(signature, Debug)
					fmt.Println()
				}
			} else {
				msg := packageName + " package does not exist. Please use --packages=true to list available packages."
				TerminalColor(msg, Error)
			}
		} else {
			header := fmt.Sprintf("\nFunctions in the %s package:", fmt.Sprintf("'%s'", lastPart))
			TerminalColor(header, Info)
			for _, groupName := range groupNames {
				functions := groupedFunctions[groupName]
				fmt.Println()
				TerminalColor("[ "+groupName+" ]", Notice)
				for _, funcInfo := range functions {
					TerminalColor("  "+strings.TrimSpace(funcInfo.Comments), Info)
					signature := fmt.Sprintf("  %s(%s) %s", funcInfo.Function, funcInfo.Params, funcInfo.Returns)
					TerminalColor(signature, Debug)
					fmt.Println()
				}
			}
		}
	}
}

// Commentify takes a string and returns it as a commented block of text,
// without adding comment syntax to the final newline if it exists.
func Commentify(str string) string {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue // Skip the empty last line
		}
		lines[i] = "  // " + line
	}
	return strings.Join(lines, "\n")
}

// Define a global map to store valid fields
var ValidFieldsMap = make(map[string]bool)

// init function to populate the ValidFieldsMap
func init() {
	ValidFields := []string{
		"CPU", "VendorID", "Family", "Model", "Stepping", "PhysicalID", "CoreID", "Cores", "ModelName",
		"Mhz", "CacheSize", "Path", "Fstype", "Total", "Free", "Used", "UsedPercent", "Available",
		"Hostname", "Uptime", "BootTime", "Procs", "OS", "Platform", "PlatformFamily", "PlatformVersion",
		"KernelVersion", "KernelArch", "VirtualizationSystem", "VirtualizationRole", "HostID",
		"Index", "Name", "HardwareAddr", "Addrs"}
	for _, field := range ValidFields {
		ValidFieldsMap[field] = true
	}
}
