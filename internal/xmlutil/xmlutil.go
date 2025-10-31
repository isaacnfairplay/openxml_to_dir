package xmlutil

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
	"unicode/utf8"
)

type Options struct {
	Indent             string
	Encoding           string
	IncludeDeclaration bool
}

func TryPrettify(name string, data []byte, opts Options) ([]byte, bool, error) {
	if !looksLikeXML(data) {
		return data, false, nil
	}
	dec := xml.NewDecoder(bytes.NewReader(data))
	dec.Strict = true
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	if opts.Indent == "" {
		enc.Indent("", "  ")
	} else {
		enc.Indent("", opts.Indent)
	}
	if opts.IncludeDeclaration {
		encoding := opts.Encoding
		if encoding == "" {
			encoding = "utf-8"
		}
		buf.WriteString("<?xml version=\"1.0\" encoding=\"" + encoding + "\"?>\n")
	}
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return data, false, err
		}
		if err := enc.EncodeToken(tok); err != nil {
			return data, false, err
		}
	}
	if err := enc.Flush(); err != nil {
		return data, false, err
	}
	return buf.Bytes(), true, nil
}

func looksLikeXML(data []byte) bool {
	if !utf8.Valid(data) {
		return false
	}
	trimmed := strings.TrimLeftFunc(string(data), func(r rune) bool {
		return r == '\n' || r == '\r' || r == '\t' || r == ' '
	})
	if trimmed == "" {
		return false
	}
	if strings.HasPrefix(trimmed, "<?xml") {
		return true
	}
	if strings.HasPrefix(trimmed, "<") && strings.Contains(trimmed, ">") {
		return true
	}
	return false
}
