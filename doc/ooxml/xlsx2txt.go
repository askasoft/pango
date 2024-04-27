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

func ExtractTextFromXlsxFile(name string) (string, error) {
	sb := &strings.Builder{}
	lw := iox.LineWriter(sb)
	err := XlsxFileTextify(name, lw)
	return sb.String(), err
}

func XlsxFileTextify(name string, w io.Writer) error {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zr.Close()

	return xlsxTextify(&zr.Reader, w)
}

func ExtractTextFromXlsxBytes(bs []byte) (string, error) {
	return ExtractTextFromXlsxReader(bytes.NewReader(bs), int64(len(bs)))
}

func ExtractTextFromXlsxReader(r io.ReaderAt, size int64) (string, error) {
	sb := &strings.Builder{}
	lw := iox.LineWriter(sb)
	err := XlsxReaderTextify(r, size, lw)
	return sb.String(), err
}

func XlsxReaderTextify(r io.ReaderAt, size int64, w io.Writer) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}

	return xlsxTextify(zr, w)
}

func xlsxTextify(zr *zip.Reader, w io.Writer) error {
	for _, zf := range zr.File {
		if zf.Name == "xl/sharedStrings.xml" {
			fr, err := zf.Open()
			if err != nil {
				return err
			}
			defer fr.Close()

			return xlsxStringsTextify(fr, w)
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
