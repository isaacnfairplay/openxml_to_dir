# ooxmlx

ooxmlx is a CGO-free command-line tool that extracts the contents of Microsoft Office OpenXML archives (such as `.docx`, `.pptx`, and `.xlsx`) into a directory, prettifying XML documents along the way.

## Usage

```bash
ooxmlx -o output_dir document.docx
```

### Flags

- `-o` (required): Destination directory for extracted files.
- `-overwrite`: Allow extraction into a non-empty destination directory.
- `-indent`: Indentation string used when formatting XML (default: two spaces).
- `-encoding`: Encoding declared in prettified XML documents (default: `utf-8`).
- `-fix-nbsp`: Normalize non-breaking spaces to regular spaces within XML payloads.
- `-quiet`: Suppress log output.
- `-version`: Print version information.

## Development

### Requirements

- Go 1.23 or newer

### Running tests

```bash
go test ./...
```

### Building

```bash
go build ./cmd/ooxmlx
```

## License

MIT
