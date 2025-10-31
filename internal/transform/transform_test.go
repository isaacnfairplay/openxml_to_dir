package transform

import "testing"

func TestReplaceNBSP(t *testing.T) {
	input := []byte("Hello\u00A0World")
	tr := ReplaceNBSP{}
	out, err := tr.Transform("test.xml", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "Hello World" {
		t.Fatalf("expected NBSP replaced, got %q", string(out))
	}
}

func TestComposite(t *testing.T) {
	comp := Composite{Transformers: []Transformer{Nop{}, ReplaceNBSP{}}}
	out, err := comp.Transform("test.xml", []byte("A\u00A0B"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "A B" {
		t.Fatalf("unexpected output %q", string(out))
	}
}
