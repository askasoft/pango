package ooxml

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/askasoft/pango/iox"
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

func ExtractTextFromDocxFile(name string) (string, error) {
	sb := &strings.Builder{}
	lw := iox.LineWriter(sb)
	err := DocxFileTextify(name, lw)
	return sb.String(), err
}

func DocxFileTextify(name string, w io.Writer) error {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zr.Close()

	return extractStringFromDocx(&zr.Reader, w)
}

func ExtractTextFromDocxBytes(bs []byte) (string, error) {
	return ExtractTextFromDocxReader(bytes.NewReader(bs), int64(len(bs)))
}

func ExtractTextFromDocxReader(r io.ReaderAt, size int64) (string, error) {
	sb := &strings.Builder{}
	lw := iox.LineWriter(sb)
	err := DocxReaderTextify(r, size, lw)
	return sb.String(), err
}

func DocxReaderTextify(r io.ReaderAt, size int64, w io.Writer) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}

	return extractStringFromDocx(zr, w)
}

func extractStringFromDocx(zr *zip.Reader, w io.Writer) error {
	for _, zf := range zr.File {
		if zf.Name == "word/document.xml" {
			fr, err := zf.Open()
			if err != nil {
				return err
			}
			defer fr.Close()

			return docxTextify(fr, w)
		}
	}
	return nil
}

func docxTextify(r io.Reader, w io.Writer) error {
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
				if _, err := iox.WriteString(w, sb.String()); err != nil {
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
			case "br", "cr":
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
