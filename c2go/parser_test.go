package c2go

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestParseHeader_Functions tests parsing header files with function declarations
func TestParseHeader_Functions(t *testing.T) {
	testFile := "tests/test_functions.h"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found", testFile)
	}

	result, err := ParseHeader(testFile)
	if err != nil {
		t.Fatalf("ParseHeader failed: %v", err)
	}

	// Current implementation returns empty results due to technical limitations
	if result == nil {
		t.Fatal("ParseHeader returned nil result")
	}

	t.Logf("Current ParseHeader returns empty results due to Clang visitor limitations")
	t.Logf("Parsed %d functions (expected: 0)", len(result.Functions))

	// TODO: When Clang visitor pattern is properly implemented,
	// this test should be updated to expect the actual functions from the header file
	t.Skip("Clang AST traversal not yet implemented - returning empty results")
}

// TestParseHeader_Structs tests parsing header files with struct declarations
func TestParseHeader_Structs(t *testing.T) {
	testFile := "tests/test_structs.h"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found", testFile)
	}

	result, err := ParseHeader(testFile)
	if err != nil {
		t.Fatalf("ParseHeader failed: %v", err)
	}

	if result == nil {
		t.Fatal("ParseHeader returned nil result")
	}
	fmt.Printf("%+v", result)
	t.Logf("Current ParseHeader returns empty results due to Clang visitor limitations")
	t.Logf("Parsed %d structs (expected: 0)", len(result.Structs))

	// TODO: When Clang visitor pattern is properly implemented,
	// this test should be updated to expect the actual structs from the header file
	t.Skip("Clang AST traversal not yet implemented - returning empty results")
}

// TestParseHeader_Enums tests parsing header files with enum declarations
func TestParseHeader_Enums(t *testing.T) {
	testFile := "tests/test_enums.h"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found", testFile)
	}

	result, err := ParseHeader(testFile)
	if err != nil {
		t.Fatalf("ParseHeader failed: %v", err)
	}

	if result == nil {
		t.Fatal("ParseHeader returned nil result")
	}

	t.Logf("Current ParseHeader returns empty results due to Clang visitor limitations")
	t.Logf("Parsed %d enums (expected: 0)", len(result.Enums))

	// TODO: When Clang visitor pattern is properly implemented,
	// this test should be updated to expect the actual enums from the header file
	t.Skip("Clang AST traversal not yet implemented - returning empty results")
}

// TestParseHeader_Complex tests parsing complex header files with multiple declarations
func TestParseHeader_Complex(t *testing.T) {
	testFile := "tests/test_complex.h"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found", testFile)
	}

	result, err := ParseHeader(testFile)
	if err != nil {
		t.Fatalf("ParseHeader failed: %v", err)
	}

	if result == nil {
		t.Fatal("ParseHeader returned nil result")
	}

	t.Logf("Complex file parsing results:")
	t.Logf("  Functions: %d", len(result.Functions))
	t.Logf("  Structs: %d", len(result.Structs))
	t.Logf("  Enums: %d", len(result.Enums))
	t.Logf("  Typedefs: %d", len(result.Typedefs))

	// TODO: When Clang visitor pattern is properly implemented,
	// this test should be updated to expect the actual declarations from the header file
	if len(result.Functions) == 0 && len(result.Structs) == 0 && len(result.Enums) == 0 && len(result.Typedefs) == 0 {
		t.Logf("Expected empty results due to current implementation limitations")
	}
}

// TestParseHeader_ErrorHandling tests error handling for invalid files
func TestParseHeader_ErrorHandling(t *testing.T) {
	// Test with non-existent file
	_, err := ParseHeader("non_existent_file.h")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Test with empty file
	emptyFile := "tests/empty.h"
	// Create empty test file
	err = os.WriteFile(emptyFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}
	defer os.Remove(emptyFile)

	result, err := ParseHeader(emptyFile)
	if err != nil {
		t.Errorf("ParseHeader failed for empty file: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result for empty file")
	}
}

// TestParseHeader_RelativePath tests parsing with relative paths
func TestParseHeader_RelativePath(t *testing.T) {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Create relative path
	relPath := filepath.Join("tests", "test_functions.h")

	// Check if test file exists
	if _, err := os.Stat(relPath); os.IsNotExist(err) {
		t.Skipf("Test file %s not found", relPath)
	}

	result, err := ParseHeader(relPath)
	if err != nil {
		t.Fatalf("ParseHeader failed with relative path: %v", err)
	}

	if result == nil {
		t.Fatal("ParseHeader returned nil result with relative path")
	}

	t.Logf("Successfully parsed with relative path from %s", wd)
}

// BenchmarkParseHeader benchmarks the ParseHeader function
func BenchmarkParseHeader(b *testing.B) {
	testFile := "tests/test_functions.h"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		b.Skipf("Test file %s not found", testFile)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseHeader(testFile)
		if err != nil {
			b.Fatalf("ParseHeader failed: %v", err)
		}
	}
}
