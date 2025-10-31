package extract

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/example/ooxmlx/internal/pathsafe"
	"github.com/example/ooxmlx/internal/transform"
	"github.com/example/ooxmlx/internal/xmlutil"
	"github.com/example/ooxmlx/internal/zipwrap"
)

type Logger interface {
	Printf(format string, v ...any)
}

type Options struct {
	Destination string
	Overwrite   bool
	Indent      string
	Encoding    string
	Transformer transform.Transformer
	Logger      Logger
}

func Extract(ctx context.Context, archivePath string, opts Options) error {
	if opts.Destination == "" {
		return errors.New("destination is required")
	}
	if opts.Transformer == nil {
		opts.Transformer = transform.Nop{}
	}
	if err := prepareDestination(opts.Destination, opts.Overwrite); err != nil {
		return err
	}
	reader, err := zipwrap.Open(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	files := reader.Files()
	for _, file := range files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := processFile(file, opts); err != nil {
			return fmt.Errorf("extract %s: %w", file.Name, err)
		}
	}
	return nil
}

func prepareDestination(dest string, overwrite bool) error {
	info, err := os.Stat(dest)
	if errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(dest, 0o755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("destination %s is not a directory", dest)
	}
	if !overwrite {
		entries, err := os.ReadDir(dest)
		if err != nil {
			return err
		}
		if len(entries) > 0 {
			return fmt.Errorf("destination %s is not empty; use overwrite", dest)
		}
	}
	return nil
}

func processFile(file zipwrap.File, opts Options) error {
	safePath, err := pathsafe.Join(opts.Destination, file.Name)
	if err != nil {
		return err
	}
	if file.IsDir {
		return os.MkdirAll(safePath, 0o755)
	}
	if err := os.MkdirAll(filepath.Dir(safePath), 0o755); err != nil {
		return err
	}
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}
	transformed, err := opts.Transformer.Transform(file.Name, data)
	if err != nil {
		return err
	}
	prettified, ok, err := xmlutil.TryPrettify(file.Name, transformed, xmlutil.Options{
		Indent:             opts.Indent,
		Encoding:           opts.Encoding,
		IncludeDeclaration: hasXMLDeclaration(transformed),
	})
	if err != nil {
		return err
	}
	payload := transformed
	if ok {
		payload = prettified
	}
	if err := os.WriteFile(safePath, payload, 0o644); err != nil {
		return err
	}
	if opts.Logger != nil {
		opts.Logger.Printf("wrote %s (%d bytes)", safePath, len(payload))
	}
	return nil
}

func hasXMLDeclaration(data []byte) bool {
	trimmed := strings.TrimLeftFunc(string(data), func(r rune) bool {
		return r == '\n' || r == '\r' || r == '\t' || r == ' '
	})
	return strings.HasPrefix(trimmed, "<?xml")
}
