package extract

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/ooxmlx/internal/transform"
)

func TestExtractSuccess(t *testing.T) {
	tmp := t.TempDir()
	archivePath := filepath.Join(tmp, "sample.zip")
	createArchive(t, archivePath, map[string][]byte{
		"doc.xml":           []byte("<root><child>value</child></root>"),
		"folder/note.txt":   []byte("plain text"),
		"folder/nested.xml": []byte("<a><b/></a>"),
	})

	dest := filepath.Join(tmp, "out")
	ctx := context.Background()
	err := Extract(ctx, archivePath, Options{
		Destination: dest,
		Overwrite:   true,
		Indent:      "  ",
		Encoding:    "utf-8",
		Transformer: transform.Nop{},
		Logger:      nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dest, "doc.xml"))
	if err != nil {
		t.Fatalf("expected doc.xml: %v", err)
	}
	if !bytes.Contains(data, []byte("<child>value</child>")) {
		t.Fatalf("unexpected xml content: %s", string(data))
	}
}

func TestExtractPathTraversal(t *testing.T) {
	tmp := t.TempDir()
	archivePath := filepath.Join(tmp, "bad.zip")
	createArchive(t, archivePath, map[string][]byte{
		"../evil.txt": []byte("nope"),
	})

	dest := filepath.Join(tmp, "out")
	ctx := context.Background()
	err := Extract(ctx, archivePath, Options{
		Destination: dest,
		Overwrite:   true,
		Indent:      "  ",
		Encoding:    "utf-8",
		Transformer: transform.Nop{},
		Logger:      nil,
	})
	if err == nil {
		t.Fatal("expected error for path traversal")
	}
}

func createArchive(t *testing.T, path string, files map[string][]byte) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create archive: %v", err)
	}
	defer file.Close()

	zw := zip.NewWriter(file)
	for name, data := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("create entry %s: %v", name, err)
		}
		if _, err := w.Write(data); err != nil {
			t.Fatalf("write entry %s: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close archive: %v", err)
	}
}
