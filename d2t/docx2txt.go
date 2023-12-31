package d2t

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

// word/document.xml
//
// ### normal text paragraph
// TEXT: This is a simple test.
//
// ```xml
// <w:p>
// 	<w:r>
// 		<w:t xml:space="preserve">This </w:t>
// 	</w:r>
// 	<w:r>
// 		<w:t>is</w:t>
// 	</w:r>
// 	<w:ins>         <!-- a inserted word -->
// 		<w:r>
// 			<w:t xml:space="preserve"> a </w:t>
// 		</w:r>
// 	</w:ins>
// 	<w:del>         <!-- a deleted history word -->
// 		<w:r>
// 			<w:delText xml:space="preserve">delete me</w:delText>
// 		</w:r>
// 	</w:del>
// 	<w:r>
// 		<w:t xml:space="preserve">simple </w:t>
// 	</w:r>
// 	<w:r>
// 		<w:t>test</w:t>
// 	</w:r>
// 	<w:r>
// 		<w:t>.</w:t>
// 	</w:r>
// </w:p>
// ```

// ### hyperlink
// ```xml
// <w:hyperlink w:history="1" w:tooltip="Adrian Boult" r:id="rId8">
// 	<w:r>
// 		<w:rPr>
// 			<w:rStyle w:val="a8"/>
// 			<w:rFonts w:cs="Arial" w:hAnsi="Arial" w:ascii="Arial"/>
// 			<w:b/>
// 			<w:bCs/>
// 			<w:color w:val="0B0080"/>
// 			<w:sz w:val="19"/>
// 			<w:szCs w:val="19"/>
// 			<w:shd w:val="clear" w:fill="F5FFFA" w:color="auto"/>
// 		</w:rPr>
// 		<w:t>Adrian Boult</w:t>
// 	</w:r>
// </w:hyperlink>
// ```

// ### tab
// ```xml
// <w:p w:rsidR="000716D9" w:rsidRDefault="000716D9">
// 	<w:r>
// 		<w:t>Banana</w:t>
// 	</w:r>
// 	<w:r>
// 		<w:tab/>
// 	</w:r>
// 	<w:r>
// 		<w:tab/>
// 	</w:r>
// 	<w:r>
// 		<w:t>banana</w:t>
// 	</w:r>
// </w:p>
// ```

// ### shift-enter (break sentence)
// ```xml
// <w:p>
// 	<w:r>
// 		<w:t>This is cat.</w:t>
// 	</w:r>
// 	<w:r>
// 		<w:br/>
// 	</w:r>
// 	<w:r>
// 		<w:t>That is dog.</w:t>
// 	</w:r>
// </w:p>
// ```

func ExtractTextFromDocxFile(name string, w io.StringWriter) error {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zr.Close()

	return extractTextFromDocx(&zr.Reader, w)
}

func ExtractTextFromDocxReader(r io.ReaderAt, size int64, w io.StringWriter) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}

	return extractTextFromDocx(zr, w)
}

func extractTextFromDocx(zr *zip.Reader, w io.StringWriter) error {
	for _, zf := range zr.File {
		if zf.Name == "word/document.xml" {
			fr, err := zf.Open()
			if err != nil {
				return err
			}
			defer fr.Close()

			return extractTextFromWordDocument(fr, w)
		}
	}
	return nil
}

func extractTextFromWordDocument(r io.Reader, w io.StringWriter) error {
	xd := xml.NewDecoder(r)

	sb := &strings.Builder{}

	wp, wt, wd := false, false, false
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
			case "del":
				wd = true
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
			case "del":
				wd = false
			case "t":
				wt = false
			case "tab":
				if wp && !wd {
					sb.WriteRune('\t')
				}
			case "br":
				if wp && !wd {
					sb.WriteRune('\n')
				}
			}
		case xml.CharData:
			if wp && wt && !wd {
				sb.Write(ty)
			}
		}
	}

	return nil
}