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

// DocxFileTextifyString Extract docx file to string
func DocxFileTextifyString(name string) (string, error) {
	sb := &strings.Builder{}
	err := DocxFileTextify(sb, name)
	return sb.String(), err
}

// DocxFileTextify Extract docx file to writer
func DocxFileTextify(w io.Writer, name string) error {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zr.Close()

	return DocxZipReaderTextify(w, &zr.Reader)
}

// DocxBytesTextifyString Extract docx data to string
func DocxBytesTextifyString(bs []byte) (string, error) {
	return DocxReaderTextifyString(bytes.NewReader(bs), int64(len(bs)))
}

// DocxBytesTextify Extract docx data to writer
func DocxBytesTextify(w io.Writer, bs []byte) error {
	return DocxReaderTextify(w, bytes.NewReader(bs), int64(len(bs)))
}

// DocxReaderTextifyString Extract docx reader to string
func DocxReaderTextifyString(r io.ReaderAt, size int64) (string, error) {
	sb := &strings.Builder{}
	err := DocxReaderTextify(sb, r, size)
	return sb.String(), err
}

// DocxReaderTextify Extract docx reader to writer
func DocxReaderTextify(w io.Writer, r io.ReaderAt, size int64) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}

	return DocxZipReaderTextify(w, zr)
}

// DocxZipReaderTextifyString Extract docx zip reader to string
func DocxZipReaderTextifyString(zr *zip.Reader) (string, error) {
	sb := &strings.Builder{}
	err := DocxZipReaderTextify(sb, zr)
	return sb.String(), err
}

// DocxZipReaderTextify Extract docx zip reader to writer
func DocxZipReaderTextify(w io.Writer, zr *zip.Reader) error {
	lw := iox.WrapWriter(w, "", "\n")

	for _, zf := range zr.File {
		if zf.Name == "word/document.xml" {
			fr, err := zf.Open()
			if err != nil {
				return err
			}
			defer fr.Close()

			return docxTextify(lw, fr)
		}
	}
	return nil
}

func docxTextify(w io.Writer, r io.Reader) error {
	var sb strings.Builder

	xd := xml.NewDecoder(r)

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
