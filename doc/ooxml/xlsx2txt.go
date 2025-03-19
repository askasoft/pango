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

// xl/sharedStrings.xml
//   TEXT: This is cell 1.
//         This is cell 2.
//
// ```xml
// <sst>
// 	<si>
// 		<r>
// 			<rPr>
// 				<b/>
// 				<sz val="36"/>
// 				<color theme="3"/>
// 				<rFont val="ＭＳ Ｐゴシック"/>
// 				<family val="3"/>
// 				<charset val="128"/>
// 				<scheme val="minor"/>
// 			</rPr>
// 			<t>This</t>
// 		</r>
// 		<r>
// 			<t xml:space="preserve"> is cell </t>
// 		</r>
// 		<r>
// 			<t>1</t>
// 		</r>
// 		<r>
// 			<t>.</t>
// 		</r>
// 	</si>
// 	<si>
// 		<r>
// 			<t xml:space="preserve">This is </t>
// 		</r>
// 		<r>
// 			<t>cell</t>
// 		</r>
// 		<r>
// 			<t xml:space="preserve"> </t>
// 		</r>
// 		<r>
// 			<t>2</t>
// 		</r>
// 		<r>
// 			<t>.</t>
// 		</r>
// 	</si>
// 	<si>
// 		<t>犬</t>
// 		<rPh sb="0" eb="1">
// 			<t>イヌ</t>
// 		</rPh>
// 		<phoneticPr fontId="1"/>
// 	</si>
// </sst>
// ```

// XlsxFileTextifyString Extract xlsx file to string
func XlsxFileTextifyString(name string) (string, error) {
	sb := &strings.Builder{}
	err := XlsxFileTextify(sb, name)
	return sb.String(), err
}

// XlsxFileTextify Extract xlsx file to writer
func XlsxFileTextify(w io.Writer, name string) error {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zr.Close()

	return XlsxZipReaderTextify(w, &zr.Reader)
}

// XlsxBytesTextifyString Extract xlsx data to string
func XlsxBytesTextifyString(bs []byte) (string, error) {
	return XlsxReaderTextifyString(bytes.NewReader(bs), int64(len(bs)))
}

// XlsxBytesTextify Extract xlsx data to writer
func XlsxBytesTextify(w io.Writer, bs []byte) error {
	return XlsxReaderTextify(w, bytes.NewReader(bs), int64(len(bs)))
}

// XlsxReaderTextifyString Extract xlsx reader to string
func XlsxReaderTextifyString(r io.ReaderAt, size int64) (string, error) {
	sb := &strings.Builder{}
	err := XlsxReaderTextify(sb, r, size)
	return sb.String(), err
}

// XlsxReaderTextify Extract xlsx reader to writer
func XlsxReaderTextify(w io.Writer, r io.ReaderAt, size int64) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}

	return XlsxZipReaderTextify(w, zr)
}

// XlsxZipReaderTextifyString Extract xlsx zip reader to string
func XlsxZipReaderTextifyString(zr *zip.Reader) (string, error) {
	sb := &strings.Builder{}
	err := XlsxZipReaderTextify(sb, zr)
	return sb.String(), err
}

// XlsxZipReaderTextify Extract xlsx zip reader to writer
func XlsxZipReaderTextify(w io.Writer, zr *zip.Reader) error {
	lw := iox.WrapWriter(w, "", "\n")

	for _, zf := range zr.File {
		if zf.Name == "xl/sharedStrings.xml" {
			fr, err := zf.Open()
			if err != nil {
				return err
			}
			defer fr.Close()

			return xlsxStringsTextify(fr, lw)
		}
	}
	return nil
}

func xlsxStringsTextify(r io.Reader, w io.Writer) error {
	var sb strings.Builder

	xd := xml.NewDecoder(r)

	xc, xt, xrph := false, false, false
	for {
		tok, err := xd.Token()
		if tok == nil || errors.Is(err, io.EOF) {
			break
		}

		switch ty := tok.(type) {
		case xml.StartElement:
			switch ty.Name.Local {
			case "si":
				xc = true
			case "rPh":
				xrph = true
			case "t":
				xt = true
			}
		case xml.EndElement:
			switch ty.Name.Local {
			case "si":
				if _, err := iox.WriteString(w, sb.String()); err != nil {
					return err
				}
				sb.Reset()
				xc = false
			case "rPh":
				xrph = false
			case "t":
				xt = false
			}
		case xml.CharData:
			if xc && xt && !xrph {
				sb.Write(ty)
			}
		}
	}

	return nil
}
