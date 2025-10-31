package pathsafe

import (
	"path/filepath"
	"testing"
)

func TestJoinInside(t *testing.T) {
	base := t.TempDir()
	path, err := Join(base, filepath.Join("foo", "bar.xml"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	want := filepath.Join(base, "foo", "bar.xml")
	if path != want {
		t.Fatalf("expected %s, got %s", want, path)
	}
}

func TestJoinOutside(t *testing.T) {
	base := t.TempDir()
	_, err := Join(base, filepath.Join("..", "evil.txt"))
	if err == nil {
		t.Fatal("expected error for escaping path")
	}
}
