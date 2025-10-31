package xmlutil

import (
	"strings"
	"testing"
)

func TestTryPrettifyXML(t *testing.T) {
	input := "<root><child>value</child></root>"
	output, ok, err := TryPrettify("test.xml", []byte(input), Options{Indent: "  ", Encoding: "utf-8", IncludeDeclaration: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected prettify to run")
	}
	result := string(output)
	if !strings.HasPrefix(result, "<?xml version=\"1.0\" encoding=\"utf-8\"?>") {
		t.Fatalf("expected xml declaration, got %s", result)
	}
	if !strings.Contains(result, "<child>value</child>") {
		t.Fatalf("unexpected output: %s", result)
	}
}

func TestTryPrettifyNonXML(t *testing.T) {
	data := []byte("plain text")
	output, ok, err := TryPrettify("text.txt", data, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected no prettify for non-xml")
	}
	if string(output) != string(data) {
		t.Fatalf("expected output to equal input")
	}
}
