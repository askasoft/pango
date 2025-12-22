package ooxml

import (
	"archive/zip"
	"bytes"
	"cmp"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/cog/treemap"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
)

// ppt/slides/slide%d.xml
//   TEXT: This is a simple test.
//
// ```xml
// <p:sld>
// 	<p:cSld>
// 		<p:spTree>
// 			<p:sp>
// 				<p:txBody>
// 					<a:p>
// 						<a:r>
// 							<a:t>This</a:t>
// 						</a:r>
// 						<a:r>
// 							<a:t> </a:t>
// 						</a:r>
// 						<a:r>
// 							<a:t>is a s</a:t>
// 						</a:r>
// 						<a:r>
// 							<a:t>imp</a:t>
// 						</a:r>
// 						<a:r>
// 							<a:t>le </a:t>
// 						</a:r>
// 						<a:r>
// 							<a:t>test</a:t>
// 						</a:r>
// 						<a:r>
// 							<a:t>.</a:t>
// 						</a:r>
// 					</a:p>
// 				</p:txBody>
// 			</p:sp>
// 		</p:spTree>
// 	</p:cSld>
// </p:sld>
// ```

// PptxFileTextifyString Extract pptx file to string
// options:
//
//	-nopgbrk             : don't insert page breaks '\f' between pages
func PptxFileTextifyString(name string, options ...string) (string, error) {
	sb := &strings.Builder{}
	err := PptxFileTextify(sb, name, options...)
	return sb.String(), err
}

// PptxFileTextify Extract pptx file to writer
// options:
//
//	-nopgbrk             : don't insert page breaks '\f' between pages
func PptxFileTextify(w io.Writer, name string, options ...string) error {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zr.Close()

	return PptxZipReaderTextify(w, &zr.Reader, options...)
}

// PptxBytesTextifyString Extract pptx data to string
// options:
//
//	-nopgbrk             : don't insert page breaks '\f' between pages
func PptxBytesTextifyString(bs []byte, options ...string) (string, error) {
	return PptxReaderTextifyString(bytes.NewReader(bs), int64(len(bs)), options...)
}

// PptxBytesTextify Extract pptx data to writer
// options:
//
//	-nopgbrk             : don't insert page breaks '\f' between pages
func PptxBytesTextify(w io.Writer, bs []byte, options ...string) error {
	return PptxReaderTextify(w, bytes.NewReader(bs), int64(len(bs)), options...)
}

// PptxReaderTextifyString Extract pptx reader to string
// options:
//
//	-nopgbrk             : don't insert page breaks '\f' between pages
func PptxReaderTextifyString(r io.ReaderAt, size int64, options ...string) (string, error) {
	sb := &strings.Builder{}
	err := PptxReaderTextify(sb, r, size, options...)
	return sb.String(), err
}

// PptxReaderTextify Extract pptx reader to writer
// options:
//
//	-nopgbrk             : don't insert page breaks '\f' between pages
func PptxReaderTextify(w io.Writer, r io.ReaderAt, size int64, options ...string) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}

	return PptxZipReaderTextify(w, zr, options...)
}

// PptxZipReaderTextifyString Extract pptx zip reader to string
// options:
//
//	-nopgbrk             : don't insert page breaks '\f' between pages
func PptxZipReaderTextifyString(zr *zip.Reader, options ...string) (string, error) {
	sb := &strings.Builder{}
	err := PptxZipReaderTextify(sb, zr, options...)
	return sb.String(), err
}

// PptxZipReaderTextify Extract pptx zip reader to writer
// options:
//
//	-nopgbrk             : don't insert page breaks '\f' between pages
func PptxZipReaderTextify(w io.Writer, zr *zip.Reader, options ...string) error {
	lw := iox.WrapWriter(w, "", "\n")

	nopgbrk := asg.Contains(options, "-nopgbrk")

	zfm := treemap.NewTreeMap[int, *zip.File](cmp.Compare[int])
	for _, zf := range zr.File {
		if str.StartsWith(zf.Name, "ppt/slides/slide") && str.EndsWith(zf.Name, ".xml") {
			zn := zf.Name[len("ppt/slides/slide") : len(zf.Name)-len(".xml")]
			if p, err := strconv.Atoi(zn); err == nil {
				zfm.Set(p, zf)
			}
		}
	}

	for it := zfm.Iterator(); it.Next(); {
		if err := pptxSlideTextify(lw, it.Value()); err != nil {
			return err
		}

		if !nopgbrk {
			if _, err := lw.Write([]byte{'\f'}); err != nil {
				return err
			}
		}
	}
	return nil
}

func pptxSlideTextify(w io.Writer, zf *zip.File) error {
	r, err := zf.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	var sb strings.Builder

	wp, wt := false, false

	for xd := xml.NewDecoder(r); ; {
		tok, err := xd.Token()
		if tok == nil || errors.Is(err, io.EOF) {
			return nil
		}

		switch ty := tok.(type) {
		case xml.StartElement:
			switch ty.Name.Local {
			case "p":
				wp = true
			case "t":
				wt = true
			}
		case xml.EndElement:
			switch ty.Name.Local {
			case "p":
				if _, err := iox.WriteString(w, sb.String()); err != nil {
					return err
				}
				sb.Reset()
				wp = false
			case "t":
				wt = false
			case "br":
				if wp {
					sb.WriteRune('\n')
				}
			}
		case xml.CharData:
			if wp && wt {
				sb.Write(ty)
			}
		}
	}
}
