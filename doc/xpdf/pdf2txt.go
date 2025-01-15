package xpdf

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

// sudo apt install poppler-utils
// Usage: pdftotext [options] <PDF-file> [<text-file>]
//   -f <int>             : first page to convert
//   -l <int>             : last page to convert
//   -r <fp>              : resolution, in DPI (default is 72)
//   -x <int>             : x-coordinate of the crop area top left corner
//   -y <int>             : y-coordinate of the crop area top left corner
//   -W <int>             : width of crop area in pixels (default is 0)
//   -H <int>             : height of crop area in pixels (default is 0)
//   -layout              : maintain original physical layout
//   -fixed <fp>          : assume fixed-pitch (or tabular) text
//   -raw                 : keep strings in content stream order
//   -nodiag              : discard diagonal text
//   -htmlmeta            : generate a simple HTML file, including the meta information
//   -tsv                 : generate a simple TSV file, including the meta information for bounding boxes
//   -enc <string>        : output text encoding name
//   -listenc             : list available encodings
//   -eol <string>        : output end-of-line convention (unix, dos, or mac)
//   -nopgbrk             : don't insert page breaks between pages
//   -bbox                : output bounding box for each word and page size to html. Sets -htmlmeta
//   -bbox-layout         : like -bbox but with extra layout bounding box data.  Sets -htmlmeta
//   -cropbox             : use the crop box rather than media box
//   -colspacing <fp>     : how much spacing we allow after a word before considering adjacent text to be a new column, as a fraction of the font size (default is 0.7, old releases had a 0.3 default)
//   -opw <string>        : owner password (for encrypted files)
//   -upw <string>        : user password (for encrypted files)
//   -q                   : don't print any messages or errors
//   -v                   : print copyright and version info
//   -h                   : print usage information
//   -help                : print usage information
//   --help               : print usage information
//   -?                   : print usage information

func PdfFileTextifyString(ctx context.Context, name string, opts ...string) (string, error) {
	bw := &bytes.Buffer{}
	err := PdfFileTextify(ctx, bw, name)
	return bw.String(), err
}

func PdfFileTextify(ctx context.Context, w io.Writer, name string, opts ...string) error {
	args := buildPdfToTextArgs(name, opts...)
	cmd := exec.CommandContext(ctx, "pdftotext", args...)
	cmd.Stdout = w
	return cmd.Run()
}

func PdfBytesTextifyString(ctx context.Context, bs []byte, opts ...string) (string, error) {
	bw := &bytes.Buffer{}
	err := PdfBytesTextify(ctx, bw, bs, opts...)
	return bw.String(), err
}

func PdfBytesTextify(ctx context.Context, w io.Writer, bs []byte, opts ...string) error {
	return PdfReaderTextify(ctx, w, bytes.NewReader(bs), opts...)
}

func PdfReaderTextifyString(ctx context.Context, r io.Reader, opts ...string) (string, error) {
	bw := &bytes.Buffer{}
	err := PdfReaderTextify(ctx, bw, r)
	return bw.String(), err
}

func PdfReaderTextify(ctx context.Context, w io.Writer, r io.Reader, opts ...string) error {
	args := buildPdfToTextArgs("-", opts...)
	cmd := exec.CommandContext(ctx, "pdftotext", args...)
	cmd.Stdin = r
	cmd.Stdout = w
	return cmd.Run()
}

func buildPdfToTextArgs(input string, opts ...string) []string {
	if len(opts) == 0 {
		opts = []string{"-layout"}
	}
	return append(opts, input, "-")
}
