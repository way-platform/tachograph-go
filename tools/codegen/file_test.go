package codegen_test

import (
	"strings"
	"testing"

	"github.com/way-platform/tacho-go/internal/codegen"
)

func TestNewFile(t *testing.T) {
	file := codegen.NewFile("test.go", "example.com/test")
	if file == nil {
		t.Fatal("NewFile should not return nil")
	}
}

func TestFile_P(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*codegen.File)
		input    []any
		expected string
	}{
		{
			name: "single string",
			setup: func(f *codegen.File) {
				// Valid Go file needs package declaration
			},
			input:    []any{"package main"},
			expected: "package main",
		},
		{
			name: "multiple strings",
			setup: func(f *codegen.File) {
				f.P("package main")
			},
			input:    []any{"func ", "main", "()", " {", "}"},
			expected: "func main() {}",
		},
		{
			name: "mixed types",
			setup: func(f *codegen.File) {
				f.P("package main")
			},
			input:    []any{"const x = ", 42},
			expected: "const x = 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := codegen.NewFile("test.go", "example.com/test")
			if tt.setup != nil {
				tt.setup(file)
			}
			file.P(tt.input...)

			content, err := file.Content()
			if err != nil {
				t.Fatalf("Content() should not fail: %v", err)
			}

			if !strings.Contains(string(content), tt.expected) {
				t.Errorf("Expected content to contain %q, but got %q", tt.expected, string(content))
			}
		})
	}
}

func TestFile_P_Empty(t *testing.T) {
	// Test P method with empty content using non-Go file
	file := codegen.NewFile("test.txt", "")
	file.P()

	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail: %v", err)
	}

	expected := "\n"
	if string(content) != expected {
		t.Errorf("Expected content %q, got %q", expected, string(content))
	}
}

func TestFile_Import(t *testing.T) {
	file := codegen.NewFile("test.go", "example.com/test")
	file.P("package main")
	file.P("func main() {}")

	// Add imports
	file.Import("fmt")
	file.Import("strings")
	file.Import("fmt") // duplicate should be ignored

	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail: %v", err)
	}

	contentStr := string(content)

	// Check that imports are present and properly formatted
	if !strings.Contains(contentStr, `import (`) {
		t.Error("Expected import block to be present")
	}
	if !strings.Contains(contentStr, `"fmt"`) {
		t.Error("Expected fmt import to be present")
	}
	if !strings.Contains(contentStr, `"strings"`) {
		t.Error("Expected strings import to be present")
	}

	// Check that imports appear before the main function
	importIndex := strings.Index(contentStr, "import")
	funcIndex := strings.Index(contentStr, "func main")
	if importIndex == -1 || funcIndex == -1 || importIndex > funcIndex {
		t.Error("Expected imports to appear before function declarations")
	}
}

func TestFile_ImportSorting(t *testing.T) {
	file := codegen.NewFile("test.go", "example.com/test")
	file.P("package main")
	file.P("func main() {}")

	// Add imports in non-alphabetical order
	file.Import("strings")
	file.Import("fmt")
	file.Import("os")
	file.Import("bufio")

	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail: %v", err)
	}

	contentStr := string(content)

	// Find positions of imports to verify sorting
	bufioPos := strings.Index(contentStr, `"bufio"`)
	fmtPos := strings.Index(contentStr, `"fmt"`)
	osPos := strings.Index(contentStr, `"os"`)
	stringsPos := strings.Index(contentStr, `"strings"`)

	if bufioPos == -1 || fmtPos == -1 || osPos == -1 || stringsPos == -1 {
		t.Fatal("All imports should be present")
	}

	// Check alphabetical order
	if !(bufioPos < fmtPos && fmtPos < osPos && osPos < stringsPos) {
		t.Errorf("Imports should be sorted alphabetically. Order: bufio=%d, fmt=%d, os=%d, strings=%d",
			bufioPos, fmtPos, osPos, stringsPos)
	}
}

func TestFile_Write(t *testing.T) {
	// Test with non-Go file to avoid parsing issues
	file := codegen.NewFile("test.txt", "")

	testData := []byte("test content")
	n, err := file.Write(testData)
	if err != nil {
		t.Fatalf("Write() should not fail: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Write() should return %d, got %d", len(testData), n)
	}

	// Write more data
	moreData := []byte(" more")
	n2, err := file.Write(moreData)
	if err != nil {
		t.Fatalf("Second Write() should not fail: %v", err)
	}
	if n2 != len(moreData) {
		t.Errorf("Second Write() should return %d, got %d", len(moreData), n2)
	}

	// Check content
	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail: %v", err)
	}

	expected := "test content more"
	if string(content) != expected {
		t.Errorf("Expected content %q, got %q", expected, string(content))
	}
}

func TestFile_ContentGoFile(t *testing.T) {
	file := codegen.NewFile("test.go", "example.com/test")
	file.Import("fmt")
	file.P("package main")
	file.P("")
	file.P("func main() {")
	file.P("    fmt.Println(\"Hello, World!\")")
	file.P("}")

	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail: %v", err)
	}

	contentStr := string(content)

	// Check that the content is properly formatted Go code
	expectedParts := []string{
		"package main",
		"import",
		`"fmt"`,
		"func main()",
		"fmt.Println",
		"Hello, World!",
	}

	for _, part := range expectedParts {
		if !strings.Contains(contentStr, part) {
			t.Errorf("Expected content to contain %q, but got:\n%s", part, contentStr)
		}
	}
}

func TestFile_ContentNonGoFile(t *testing.T) {
	file := codegen.NewFile("test.txt", "")
	file.P("This is a text file")
	file.P("with multiple lines")

	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail: %v", err)
	}

	expected := "This is a text file\nwith multiple lines\n"
	if string(content) != expected {
		t.Errorf("Expected content %q, got %q", expected, string(content))
	}
}

func TestFile_ContentInvalidGoCode(t *testing.T) {
	file := codegen.NewFile("test.go", "example.com/test")
	file.P("package main")
	file.P("func main() {")
	file.P("    invalid syntax here }")
	file.P("missing closing brace")

	_, err := file.Content()
	if err == nil {
		t.Error("Content() should fail for invalid Go code")
	}

	// Check that error message contains helpful information
	errorStr := err.Error()
	if !strings.Contains(errorStr, "unparsable Go source") {
		t.Errorf("Error should mention unparsable Go source, got: %v", err)
	}
	if !strings.Contains(errorStr, "test.go") {
		t.Errorf("Error should mention filename, got: %v", err)
	}
}

func TestFile_ContentWithComments(t *testing.T) {
	file := codegen.NewFile("test.go", "example.com/test")
	file.Import("fmt")
	file.P("// Package main provides a simple example")
	file.P("package main")
	file.P("")
	file.P("// main is the entry point")
	file.P("func main() {")
	file.P("    fmt.Println(\"Hello\")")
	file.P("}")

	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail: %v", err)
	}

	contentStr := string(content)

	// Check that comments are preserved
	if !strings.Contains(contentStr, "Package main provides") {
		t.Error("Package comment should be preserved")
	}
	if !strings.Contains(contentStr, "main is the entry point") {
		t.Error("Function comment should be preserved")
	}

	// Check that imports still come after package declaration but before functions
	packagePos := strings.Index(contentStr, "package main")
	importPos := strings.Index(contentStr, "import")
	funcPos := strings.Index(contentStr, "func main")

	if packagePos == -1 || importPos == -1 || funcPos == -1 {
		t.Fatal("Package, import, and function declarations should all be present")
	}

	if !(packagePos < importPos && importPos < funcPos) {
		t.Error("Expected order: package, import, function")
	}
}

func TestFile_MultipleOperations(t *testing.T) {
	file := codegen.NewFile("complex.go", "example.com/test")

	// Test io.Writer interface
	_, err := file.Write([]byte("// Generated code\n"))
	if err != nil {
		t.Fatalf("Write() should not fail: %v", err)
	}

	// Test P method
	file.P("package main")
	file.P("")

	// Test Import method
	file.Import("fmt")
	file.Import("os")

	// Add more content
	file.P("func main() {")
	file.P("    fmt.Println(\"Starting program\")")
	file.P("    if len(os.Args) > 1 {")
	file.P("        fmt.Printf(\"Args: %v\\n\", os.Args[1:])")
	file.P("    }")
	file.P("}")

	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail: %v", err)
	}

	contentStr := string(content)

	// Verify all expected parts are present
	expectedParts := []string{
		"Generated code",
		"package main",
		"import",
		`"fmt"`,
		`"os"`,
		"func main",
		"fmt.Println",
		"os.Args",
	}

	for _, part := range expectedParts {
		if !strings.Contains(contentStr, part) {
			t.Errorf("Expected content to contain %q, but content was:\n%s", part, contentStr)
		}
	}
}

func TestFile_EmptyFile(t *testing.T) {
	// Test with non-Go file since empty Go files are invalid
	file := codegen.NewFile("empty.txt", "")

	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail for empty file: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("Expected empty content, got %q", string(content))
	}
}

func TestFile_OnlyPackage(t *testing.T) {
	file := codegen.NewFile("package.go", "example.com/test")
	file.P("package test")

	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "package test") {
		t.Errorf("Expected content to contain package declaration, got %q", contentStr)
	}
}

func TestFile_ImportWithoutCode(t *testing.T) {
	file := codegen.NewFile("imports.go", "example.com/test")
	file.Import("fmt")
	file.Import("os")

	// This should fail because there's no package declaration
	_, err := file.Content()
	if err == nil {
		t.Error("Content() should fail when imports are added but no package is declared")
	}
}

func TestFile_QualifiedGoIdent(t *testing.T) {
	file := codegen.NewFile("test.go", "example.com/test")
	
	// Test same package - should return just the name
	samePackageIdent := codegen.GoIdent{GoImportPath: "example.com/test", GoName: "MyType"}
	result := file.QualifiedGoIdent(samePackageIdent)
	if result != "MyType" {
		t.Errorf("Expected 'MyType', got '%s'", result)
	}
	
	// Test built-in type - should return just the name
	builtinIdent := codegen.GoIdent{GoImportPath: "", GoName: "string"}
	result = file.QualifiedGoIdent(builtinIdent)
	if result != "string" {
		t.Errorf("Expected 'string', got '%s'", result)
	}
	
	// Test different package - should return qualified name and add import
	xmlIdent := codegen.GoIdent{GoImportPath: "encoding/xml", GoName: "Name"}
	result = file.QualifiedGoIdent(xmlIdent)
	if result != "xml.Name" {
		t.Errorf("Expected 'xml.Name', got '%s'", result)
	}
	
	// Test package name collision handling
	timeIdent := codegen.GoIdent{GoImportPath: "time", GoName: "Time"}
	result = file.QualifiedGoIdent(timeIdent)
	if result != "time.Time" {
		t.Errorf("Expected 'time.Time', got '%s'", result)
	}
	
	// Test that imports were added
	file.P("package test")
	content, err := file.Content()
	if err != nil {
		t.Fatalf("Content() should not fail: %v", err)
	}
	
	contentStr := string(content)
	if !strings.Contains(contentStr, `"encoding/xml"`) {
		t.Error("Expected xml import to be added")
	}
	if !strings.Contains(contentStr, `"time"`) {
		t.Error("Expected time import to be added")
	}
}

func TestFile_QualifiedGoIdent_PackageNameCollision(t *testing.T) {
	file := codegen.NewFile("test.go", "example.com/test")
	
	// Create two packages that would have the same base name
	pkg1Ident := codegen.GoIdent{GoImportPath: "example.com/foo/xml", GoName: "Parser"}
	pkg2Ident := codegen.GoIdent{GoImportPath: "encoding/xml", GoName: "Name"}
	
	result1 := file.QualifiedGoIdent(pkg1Ident)
	result2 := file.QualifiedGoIdent(pkg2Ident)
	
	// Both should be qualified but with different package names
	if result1 == result2 {
		t.Errorf("Expected different qualifications, got same: %s", result1)
	}
	
	// Should contain the type names
	if !strings.Contains(result1, "Parser") {
		t.Errorf("Expected result1 to contain 'Parser', got '%s'", result1)
	}
	if !strings.Contains(result2, "Name") {
		t.Errorf("Expected result2 to contain 'Name', got '%s'", result2)
	}
}

func TestFile_CommonIdents(t *testing.T) {
	file := codegen.NewFile("test.go", "example.com/test")
	
	// Test using common identifiers
	xmlNameResult := file.QualifiedGoIdent(codegen.XMLNameIdent)
	contextResult := file.QualifiedGoIdent(codegen.ContextIdent)
	stringResult := file.QualifiedGoIdent(codegen.StringIdent)
	
	if xmlNameResult != "xml.Name" {
		t.Errorf("Expected 'xml.Name', got '%s'", xmlNameResult)
	}
	if contextResult != "context.Context" {
		t.Errorf("Expected 'context.Context', got '%s'", contextResult)
	}
	if stringResult != "string" {
		t.Errorf("Expected 'string', got '%s'", stringResult)
	}
}
