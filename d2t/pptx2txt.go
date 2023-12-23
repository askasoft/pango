package d2t

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/askasoft/pango/cog"
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

func ExtractTextFromPptxFile(name string, w io.StringWriter) error {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zr.Close()

	return extractTextFromPptx(&zr.Reader, w)
}

func ExtractTextFromPptxReader(r io.ReaderAt, size int64, w io.StringWriter) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}

	return extractTextFromPptx(zr, w)
}

func extractTextFromPptx(zr *zip.Reader, w io.StringWriter) error {
	zfm := cog.NewTreeMap[string, *zip.File](cog.CompareString)
	for _, zf := range zr.File {
		if str.StartsWith(zf.Name, "ppt/slides/slide") && str.EndsWith(zf.Name, ".xml") {
			zfm.Set(zf.Name, zf)
		}
	}

	for it := zfm.Iterator(); it.Next(); {
		fr, err := it.Value().Open()
		if err != nil {
			return err
		}
		defer fr.Close()

		if err := extractTextFromPptxSlide(fr, w); err != nil {
			return err
		}
	}
	return nil
}

func extractTextFromPptxSlide(r io.Reader, w io.StringWriter) error {
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
				if _, err := w.WriteString(sb.String()); err != nil {
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
