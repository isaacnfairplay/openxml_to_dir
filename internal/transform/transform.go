package transform

import "unicode/utf8"

type Transformer interface {
	Transform(name string, data []byte) ([]byte, error)
}

type Nop struct{}

func (Nop) Transform(_ string, data []byte) ([]byte, error) {
	return data, nil
}

type Composite struct {
	Transformers []Transformer
}

func (c Composite) Transform(name string, data []byte) ([]byte, error) {
	current := data
	var err error
	for _, t := range c.Transformers {
		current, err = t.Transform(name, current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}

type ReplaceNBSP struct{}

func (ReplaceNBSP) Transform(_ string, data []byte) ([]byte, error) {
	// Fast path: check if NBSP present
	contains := false
	for i := 0; i < len(data); {
		r, size := utf8.DecodeRune(data[i:])
		if r == utf8.RuneError && size == 1 {
			i++
			continue
		}
		if r == 0x00A0 {
			contains = true
			break
		}
		i += size
	}
	if !contains {
		return data, nil
	}
	out := make([]byte, 0, len(data))
	for i := 0; i < len(data); {
		r, size := utf8.DecodeRune(data[i:])
		if r == 0x00A0 {
			out = append(out, ' ')
		} else {
			out = append(out, data[i:i+size]...)
		}
		if r == utf8.RuneError && size == 1 {
			i++
		} else {
			i += size
		}
	}
	return out, nil
}
