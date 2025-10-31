package zipwrap

import (
	"archive/zip"
	"io"
	"os"
)

type Reader struct {
	file   *os.File
	reader *zip.Reader
}

type File struct {
	Name  string
	Size  uint64
	IsDir bool
	Open  func() (io.ReadCloser, error)
}

func Open(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	zr, err := zip.NewReader(f, info.Size())
	if err != nil {
		f.Close()
		return nil, err
	}
	return &Reader{file: f, reader: zr}, nil
}

func (r *Reader) Close() error {
	if r.file == nil {
		return nil
	}
	err := r.file.Close()
	r.file = nil
	r.reader = nil
	return err
}

func (r *Reader) Files() []File {
	if r.reader == nil {
		return nil
	}
	out := make([]File, 0, len(r.reader.File))
	for _, f := range r.reader.File {
		file := f
		info := file.FileInfo()
		out = append(out, File{
			Name:  file.Name,
			Size:  file.UncompressedSize64,
			IsDir: info.IsDir(),
			Open: func() (io.ReadCloser, error) {
				if info.IsDir() {
					return io.NopCloser(&emptyReader{}), nil
				}
				return file.Open()
			},
		})
	}
	return out
}

type emptyReader struct{}

func (e *emptyReader) Read(p []byte) (int, error) {
	return 0, io.EOF
}
