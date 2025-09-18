package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"
)

// GoIdent represents a Go identifier with its import path.
type GoIdent struct {
	GoImportPath string // Import path, e.g., "encoding/xml"
	GoName       string // Identifier name, e.g., "Name"
}

// GoPackageName represents a package name used in generated code.
type GoPackageName string

// NewFile creates a new [File] with the given filename and import path.
func NewFile(filename, goImportPath string) *File {
	return &File{
		filename:         filename,
		goImportPath:     goImportPath,
		packageNames:     make(map[string]GoPackageName),
		usedPackageNames: make(map[GoPackageName]bool),
		imports:          make(map[string]bool),
	}
}

// File represents a generated Go file.
type File struct {
	filename         string
	goImportPath     string                   // Import path of the generated file
	packageNames     map[string]GoPackageName // Import path -> package name
	usedPackageNames map[GoPackageName]bool   // Track used package names
	imports          map[string]bool          // Import paths to include
	buf              bytes.Buffer
}

// P writes a line of code to the file.
func (g *File) P(v ...any) {
	for _, x := range v {
		fmt.Fprint(&g.buf, x)
	}
	fmt.Fprintln(&g.buf)
}

// QualifiedGoIdent returns the qualified Go identifier and manages imports automatically.
// If the identifier is from the same package, returns just the name.
// If from a different package, returns packageName.Name and ensures the import is added.
func (f *File) QualifiedGoIdent(ident GoIdent) string {
	if ident.GoImportPath == f.goImportPath || ident.GoImportPath == "" {
		return ident.GoName
	}

	if packageName, ok := f.packageNames[ident.GoImportPath]; ok {
		f.imports[ident.GoImportPath] = true
		return string(packageName) + "." + ident.GoName
	}

	packageName := cleanPackageName(path.Base(ident.GoImportPath))
	for i, orig := 1, packageName; f.usedPackageNames[GoPackageName(packageName)]; i++ {
		packageName = orig + strconv.Itoa(i)
	}

	f.packageNames[ident.GoImportPath] = GoPackageName(packageName)
	f.usedPackageNames[GoPackageName(packageName)] = true
	f.imports[ident.GoImportPath] = true

	return packageName + "." + ident.GoName
}

// Import adds a blank import to the file (for side effects).
// For normal imports, use QualifiedGoIdent instead.
func (f *File) Import(importPath string) {
	f.imports[importPath] = true
}

// SetPackageName sets a custom package name for the given import path.
// This allows overriding the default package name derivation.
func (f *File) SetPackageName(importPath string, packageName string) {
	f.packageNames[importPath] = GoPackageName(packageName)
	f.usedPackageNames[GoPackageName(packageName)] = true
}

// cleanPackageName returns a valid Go package name from an import path.
func cleanPackageName(name string) string {
	// Remove common suffixes and invalid characters
	name = strings.TrimSuffix(name, ".go")
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, ".", "")

	// Ensure it starts with a letter
	if len(name) == 0 || !((name[0] >= 'a' && name[0] <= 'z') || (name[0] >= 'A' && name[0] <= 'Z')) {
		name = "pkg" + name
	}

	return name
}

// Filename returns the filename of the file.
func (g *File) Filename() string {
	return g.filename
}

// Write implements [io.Writer].
func (g *File) Write(p []byte) (n int, err error) {
	return g.buf.Write(p)
}

// Content returns the contents of the generated file.
func (g *File) Content() ([]byte, error) {
	if !strings.HasSuffix(g.filename, ".go") {
		return g.buf.Bytes(), nil
	}
	// Reformat generated code.
	original := g.buf.Bytes()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", original, parser.ParseComments)
	if err != nil {
		// Print out the bad code with line numbers.
		// This should never happen in practice, but it can while changing generated code
		// so consider this a debugging aid.
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(original))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		return nil, fmt.Errorf("%v: unparsable Go source: %v\n%v", g.filename, err, src.String())
	}
	// Collect a sorted list of all imports.
	var importPaths []string
	for importPath := range g.imports {
		importPaths = append(importPaths, importPath)
	}
	sort.Strings(importPaths)
	// Modify the AST to include a new import block.
	if len(importPaths) > 0 {
		// Insert block after package statement or
		// possible comment attached to the end of the package statement.
		pos := file.Package
		tokFile := fset.File(file.Package)
		pkgLine := tokFile.Line(file.Package)
		for _, c := range file.Comments {
			if tokFile.Line(c.Pos()) > pkgLine {
				break
			}
			pos = c.End()
		}
		// Construct the import block.
		impDecl := &ast.GenDecl{
			Tok:    token.IMPORT,
			TokPos: pos,
			Lparen: pos,
			Rparen: pos,
		}
		for _, importPath := range importPaths {
			spec := &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:     token.STRING,
					Value:    strconv.Quote(importPath),
					ValuePos: pos,
				},
				EndPos: pos,
			}
			// Add alias if we have a custom package name for this import
			if packageName, ok := g.packageNames[importPath]; ok {
				expectedName := cleanPackageName(path.Base(importPath))
				if string(packageName) != expectedName {
					spec.Name = &ast.Ident{
						Name:    string(packageName),
						NamePos: pos,
					}
				}
			}
			impDecl.Specs = append(impDecl.Specs, spec)
		}
		file.Decls = append([]ast.Decl{impDecl}, file.Decls...)
	}
	var out bytes.Buffer
	if err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(&out, fset, file); err != nil {
		return nil, fmt.Errorf("%v: can not reformat Go source: %v", g.filename, err)
	}
	return out.Bytes(), nil
}
