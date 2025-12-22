package ooxml

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
)

// @see http://officeopenxml.com/WPcontentOverview.php

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

	numFile := docxFindZipFile(zr, "word/numbering.xml")
	numDefs, _ := docxParseNumbering(numFile)

	docFile := docxFindZipFile(zr, "word/document.xml")
	if docFile == nil {
		return errors.New("docx: missing word/document.xml")
	}

	return docxTextify(lw, docFile, numDefs)
}

func docxFindZipFile(zr *zip.Reader, name string) *zip.File {
	for _, f := range zr.File {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// docxTextify parses the document.xml content and writes lines to writer,
// inserting numbering prefixes for ordered list paragraphs where possible.
func docxTextify(w io.Writer, zf *zip.File, nd *numDefinitions) error {
	r, err := zf.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	var sb strings.Builder

	wp, wt, wd, np := false, false, false, ""

	for xd := xml.NewDecoder(r); ; {
		tok, err := xd.Token()
		if tok == nil || errors.Is(err, io.EOF) {
			return nil
		}

		switch ty := tok.(type) {
		case xml.StartElement:
			switch ty.Name.Local {
			case "p":
				// new paragraph
				wp, wt, wd, np = true, false, false, ""
				sb.Reset()
			case "pPr":
				// paragraph properties: try to extract numbering
				var pp pPr
				if err := xd.DecodeElement(&pp, &ty); err == nil {
					if pp.NumPr != nil && pp.NumPr.NumID != nil {
						numId := pp.NumPr.NumID.Val
						ilvl := 0
						if pp.NumPr.ILvl != nil {
							ilvl = pp.NumPr.ILvl.Val
						}
						np = nd.next(numId, ilvl)
					}
				}
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
				if sb.Len() == 0 && np != "" {
					sb.WriteString(np)
					sb.WriteByte(' ')
					np = ""
				}
				sb.Write(ty)
			}
		}
	}
}

// --- minimal numbering XML structures we need ---

type pPr struct {
	NumPr *numPr `xml:"numPr"`
}

type numPr struct {
	NumID *valInt `xml:"numId"`
	ILvl  *valInt `xml:"ilvl"`
}

// Note: we no longer model run/hyperlink/text nodes since we switched
// to a token-based parser to preserve ordering of tabs, breaks, etc.

type valInt struct {
	Val int `xml:"val,attr"`
}

type valStr struct {
	Val string `xml:"val,attr"`
}

// --- Numbering support ---

type numbering struct {
	Nums         []num         `xml:"num"`
	AbstractNums []abstractNum `xml:"abstractNum"`
}

type num struct {
	NumID         int           `xml:"numId,attr"`
	AbstractNumID *valInt       `xml:"abstractNumId"`
	Lvls          []lvlOverride `xml:"lvlOverride"`
}

type lvlOverride struct {
	Ilvl  int     `xml:"ilvl,attr"`
	Start *valInt `xml:"startOverride>start"`
}

type abstractNum struct {
	AbstractNumID int   `xml:"abstractNumId,attr"`
	Lvls          []lvl `xml:"lvl"`
}

type lvl struct {
	Ilvl    int    `xml:"ilvl,attr"`
	NumFmt  valStr `xml:"numFmt"`
	LvlText valStr `xml:"lvlText"`
	Start   valInt `xml:"start"`
}

// numDefinitions keeps counters and level formatting for each list instance.
type numDefinitions struct {
	numToAbs map[int]int
	absLvls  map[int]map[int]lvl // absId -> ilvl -> lvl
	counters map[int][]int       // numId -> counters per level
	starts   map[int]map[int]int // numId -> ilvl -> start
}

func (nd *numDefinitions) next(numId int, ilvl int) string {
	if nd == nil {
		return ""
	}

	absId, ok := nd.numToAbs[numId]
	if !ok {
		// unknown definition
		return ""
	}

	lmap := nd.absLvls[absId]
	l, ok := lmap[ilvl]
	if !ok {
		return ""
	}

	// ensure counters slice is large enough
	cs := nd.counters[numId]
	if len(cs) <= ilvl {
		ncs := make([]int, ilvl+1)
		copy(ncs, cs)
		cs = ncs
	}

	// determine start value (default 1)
	start := 1
	if s, ok := nd.starts[numId][ilvl]; ok && s > 0 {
		start = s
	} else if l.Start.Val > 0 {
		start = l.Start.Val
	}

	// increment this level, reset deeper levels
	if cs[ilvl] == 0 {
		cs[ilvl] = start
	} else {
		cs[ilvl]++
	}
	for i := ilvl + 1; i < len(cs); i++ {
		cs[i] = 0
	}
	nd.counters[numId] = cs

	// If it's an explicit bullet list, keep bullet ("-"); otherwise, render using lvlText.
	// This makes numeric formats such as decimalFullWidth, decimalEnclosedCircle, etc. display as digits.
	if str.EqualFold(l.NumFmt.Val, "bullet") {
		return "-"
	}
	return formatLvlText(l.LvlText.Val, cs[:ilvl+1])
}

func formatLvlText(tpl string, counts []int) string {
	// Replace %1, %2, ... in lvlText with corresponding counts converted as needed.
	// We only support decimal explicitly; for other formats above we fall back to decimal output.
	res := tpl

	// DOCX may use "%1." or "%1)" etc. We replace all "%n" occurrences we know about.
	for i, c := range counts {
		res = strings.ReplaceAll(res, "%"+strconv.Itoa(i+1), strconv.Itoa(c))
	}

	// If tpl did not contain placeholders, just use the deepest level number with a dot.
	if res == tpl {
		res = strconv.Itoa(counts[len(counts)-1]) + "."
	}
	return res
}

func docxParseNumbering(zf *zip.File) (*numDefinitions, error) {
	if zf == nil {
		return nil, nil
	}

	r, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var n numbering
	err = xml.NewDecoder(r).Decode(&n)
	if err != nil {
		return nil, err
	}

	nd := &numDefinitions{
		numToAbs: map[int]int{},
		absLvls:  map[int]map[int]lvl{},
		counters: map[int][]int{},
		starts:   map[int]map[int]int{},
	}

	for _, an := range n.AbstractNums {
		lmap := map[int]lvl{}

		for _, l := range an.Lvls {
			lmap[l.Ilvl] = l
		}
		nd.absLvls[an.AbstractNumID] = lmap
	}

	for _, m := range n.Nums {
		if m.AbstractNumID != nil {
			nd.numToAbs[m.NumID] = m.AbstractNumID.Val
		}
		if len(m.Lvls) > 0 {
			if nd.starts[m.NumID] == nil {
				nd.starts[m.NumID] = map[int]int{}
			}
			for _, ov := range m.Lvls {
				if ov.Start != nil {
					nd.starts[m.NumID][ov.Ilvl] = ov.Start.Val
				}
			}
		}
	}

	return nd, nil
}
