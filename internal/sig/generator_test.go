package sig

import (
	"os"
	"testing"
)

func TestParseV2(t *testing.T) {
	pkgInfo, err := Parse("./sigtest")
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
