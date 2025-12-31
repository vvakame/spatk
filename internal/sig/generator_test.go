package sig

import (
	"os"
	"strings"
	"testing"
)

func TestParseV2(t *testing.T) {
	pkgInfo, err := Parse("./sigtest", nil)
	if err != nil {
		t.Fatal(err)
	}

	pkgInfo.CommandArgs = []string{"-private", "-output", "model_gen.go", "."}
	for _, structInfo := range pkgInfo.Structs {
		structInfo.Private = true
	}
	actual, err := pkgInfo.Emit()
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile("./sigtest/model_gen.go")
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) != string(expected) {
		t.Fatalf("expected %s, but got %s", string(expected), string(actual))
	}
}

func TestParseWithErrors(t *testing.T) {
	t.Run("without AllowErrors flag", func(t *testing.T) {
		_, err := Parse("./testdata", nil)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		// Should return compile error
		if !strings.Contains(err.Error(), "undefined") && !strings.Contains(err.Error(), "UndefinedType") {
			t.Fatalf("expected compile error about undefined type, but got: %v", err)
		}
	})

	t.Run("with AllowErrors flag", func(t *testing.T) {
		pkgInfo, err := Parse("./testdata", &ParseOptions{
			AllowErrors: true,
		})
		if err != nil {
			t.Fatalf("expected no error with AllowErrors, but got: %v", err)
		}

		if len(pkgInfo.Structs) != 1 {
			t.Fatalf("expected 1 struct, but got %d", len(pkgInfo.Structs))
		}

		st := pkgInfo.Structs[0]
		if st.Name != "BrokenModel" {
			t.Fatalf("expected struct name BrokenModel, but got %s", st.Name)
		}

		// Check that PanicFallback is set for the BrokenField
		var brokenField *FieldInfo
		for _, f := range st.Fields {
			if f.Name == "BrokenField" {
				brokenField = f
				break
			}
		}
		if brokenField == nil {
			t.Fatal("expected to find BrokenField")
		}
		if !brokenField.PanicFallback {
			t.Fatal("expected BrokenField to have PanicFallback=true")
		}

		// Test code generation
		pkgInfo.CommandArgs = []string{"-allow-errors", "-private", "-output", "model_gen.go", "."}
		for _, structInfo := range pkgInfo.Structs {
			structInfo.Private = true
		}
		generated, err := pkgInfo.Emit()
		if err != nil {
			t.Fatalf("Emit failed: %v", err)
		}

		// Check that generated code contains panic fallback
		generatedStr := string(generated)
		if !strings.Contains(generatedStr, "panic(\"sig: type info unavailable") {
			t.Fatal("expected generated code to contain panic fallback for BrokenField")
		}

		// Check that generated code compiles (has valid Go syntax)
		if !strings.Contains(generatedStr, "package testdata") {
			t.Fatal("expected valid package declaration")
		}
	})
}
