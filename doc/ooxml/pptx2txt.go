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

func PptxFileTextifyString(name string, opts ...string) (string, error) {
	sb := &strings.Builder{}
	err := PptxFileTextify(sb, name, opts...)
	return sb.String(), err
}

func PptxFileTextify(w io.Writer, name string, opts ...string) error {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zr.Close()

	return PptxZipReaderTextify(w, &zr.Reader, opts...)
}

func PptxBytesTextifyString(bs []byte, opts ...string) (string, error) {
	return PptxReaderTextifyString(bytes.NewReader(bs), int64(len(bs)), opts...)
}

func PptxBytesTextify(w io.Writer, bs []byte, opts ...string) error {
	return PptxReaderTextify(w, bytes.NewReader(bs), int64(len(bs)), opts...)
}

func PptxReaderTextifyString(r io.ReaderAt, size int64, opts ...string) (string, error) {
	sb := &strings.Builder{}
	err := PptxReaderTextify(sb, r, size, opts...)
	return sb.String(), err
}

func PptxReaderTextify(w io.Writer, r io.ReaderAt, size int64, opts ...string) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}

	return PptxZipReaderTextify(w, zr, opts...)
}

func PptxZipReaderTextifyString(zr *zip.Reader, opts ...string) (string, error) {
	sb := &strings.Builder{}
	err := PptxZipReaderTextify(sb, zr, opts...)
	return sb.String(), err
}

func PptxZipReaderTextify(w io.Writer, zr *zip.Reader, opts ...string) error {
	lw := iox.WrapWriter(w, "", "\n")

	nopgbrk := asg.Contains(opts, "-nopgbrk")

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
		fr, err := it.Value().Open()
		if err != nil {
			return err
		}
		defer fr.Close()

		if err := pptxSlideTextify(lw, fr); err != nil {
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

func pptxSlideTextify(w io.Writer, r io.Reader) error {
	xd := xml.NewDecoder(r)

	sb := &strings.Builder{}

	wp, wt := false, false
	for {
		tok, err := xd.Token()
		if tok == nil || errors.Is(err, io.EOF) {
			break
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

	return nil
}
