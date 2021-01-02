package pdf

import (
	"bufio"
	"bytes"

	"github.com/jung-kurt/gofpdf"
)

type PDFTable struct {
	Headers []PDFCell
	Body    [][]PDFCell
}

type PDFCell struct {
	value     string
	align     string
	fontStyle string
}

func NewPDFCell(value, align, style string) PDFCell {
	a := align
	if align != "C" && align != "R" {
		a = "L"
	}
	return PDFCell{
		value:     value,
		align:     a,
		fontStyle: style,
	}
}

func NewPDFROW(values []string, align, style string) []PDFCell {
	row := make([]PDFCell, len(values))
	for i, v := range values {
		row[i] = NewPDFCell(v, align, style)
	}
	return row
}

func (p PDFCell) String() string {
	return p.value
}

// func PrintHeading(pdf *gofpdf.Fpdf, heading string) {

// 	pdf.SetX(0.0)
// 	_, f := pdf.GetFontSize()
// 	pdf.WriteAligned(0, f+2, heading, "C")
// 	pdf.Ln(-1)

// }

func PrintText(pdf *gofpdf.Fpdf, str string) {

	//pdf.SetX(0.0)
	_, f := pdf.GetFontSize()
	pdf.Write(f+2, str)
	//pdf.Ln(-1)

}

func PrintTable(pdf *gofpdf.Fpdf, t PDFTable, colwr, colw []float64, w float64) {
	if w == 0 {
		pageWidth, _ := pdf.GetPageSize()
		lMargin, _, rMargin, _ := pdf.GetMargins()
		w = pageWidth - lMargin - rMargin
	}
	//log.Println("Colwidth ration:", colwr)
	cw := DoColWidths(pdf, t, colwr, colw, w)
	//log.Println("Colwidth:", cw)

	fs, _ := pdf.GetFontSize()

	if t.Headers != nil {
		pdf.SetFontSize(fs + 2.0)
		PrintRow(pdf, t.Headers, cw)
		pdf.SetFontSize(fs)
	}

	for _, row := range t.Body {
		PrintRow(pdf, row, cw)
	}

}

func PrintRow(pdf *gofpdf.Fpdf, row []PDFCell, cols []float64) {
	_, pageh := pdf.GetPageSize()
	_, _, _, mbottom := pdf.GetMargins()

	_, lineHt := pdf.GetFontSize()

	curx, y := pdf.GetXY()
	x := curx
	height := 0.
	//determince cell height
	for i, txt := range row {
		lines := pdf.SplitLines([]byte(txt.value), cols[i])
		h := float64(len(lines)) * (lineHt + 2.0)
		if h > height {
			height = h
		}
	}

	// add a new page if the height of the row doesn't fit on the page
	if pdf.GetY()+height > pageh-mbottom {
		pdf.AddPage()
		y = pdf.GetY()
	}

	//write row
	for i, txt := range row {
		width := cols[i]
		pdf.Rect(x, y, width, height, "")
		pdf.SetFontStyle(txt.fontStyle)
		pdf.MultiCell(width, lineHt+2.0, txt.value, "", txt.align, false)
		pdf.SetFontStyle("")
		x += width
		pdf.SetXY(x, y)
	}
	pdf.SetXY(curx, y+height)

}

func DoColWidths(pdf *gofpdf.Fpdf, t PDFTable, colwr, colw []float64, w float64) []float64 {
	if colw != nil {
		return colw
	}

	if colwr != nil {
		return ColWidthsFromRatio(colwr, w)
	}
	return GenColWidts(pdf, t, w)
}

func ColWidthsFromRatio(colwr []float64, w float64) []float64 {
	rl := len(colwr)
	cw := make([]float64, rl)
	sr := 0.0
	for _, c := range colwr {
		sr += c
	}
	ratio := w / sr
	for i := 0; i < rl; i++ {
		cw[i] = ratio * colwr[i]
	}
	return cw
}

func GenColWidts(pdf *gofpdf.Fpdf, t PDFTable, w float64) []float64 {

	var rl int
	if t.Headers == nil {
		rl = len(t.Body[0])
	} else {
		rl = len(t.Headers)
	}

	av := make([]float64, rl)
	cw := make([]float64, rl)
	var toDefWrap bool
	maxW := make([]float64, rl)
	defwidth := w / float64(rl)

	if t.Headers != nil {
		for i, r := range t.Headers {
			k := pdf.GetStringWidth(r.value)
			av[i] += k
			maxW[i] = max(maxW[i], k)
		}
	}

	for _, row := range t.Body {

		for i, r := range row {
			k := pdf.GetStringWidth(r.value)
			av[i] += k
			maxW[i] = max(maxW[i], k)
		}
	}

	bl := float64(len(t.Body))

	for i := 0; i < rl; i++ {
		av[i] = av[i] / (bl + 1)
		if av[i] > defwidth-2 {
			toDefWrap = true
		}
	}

	if !toDefWrap {
		for i := 0; i < rl; i++ {
			cw[i] = defwidth
		}
	} else {
		//get max per col, set width for the ones
		//with widt that will fit default cell width and
		//assign the default width if max with is bigger than
		//default withd or assigh max with if less
		s := 0.0
		for i := 0; i < rl; i++ {
			if av[i] <= defwidth-2 {
				if maxW[i] <= defwidth-2 {
					cw[i] = maxW[i] + 2
				} else {
					cw[i] = defwidth
				}
				s += cw[i]
			}
		}
		//share the rest of space with the big cols
		s = w - s
		sum := 0.0
		//sum of the average witds of remaining coll
		for i := 0; i < rl; i++ {
			if cw[i] == 0 {
				sum += av[i]
			}
		}
		// share it in the ration of their av width
		ratio := s / sum
		for i := 0; i < rl; i++ {
			if cw[i] == 0 {
				cw[i] = ratio * av[i]
			}
		}

	}
	return cw
}

func GetInitPDF(addPage bool) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	if addPage {
		pdf.AddPage()
	}

	return pdf
}

func PDF2Bytes(pdf *gofpdf.Fpdf) ([]byte, error) {
	var b bytes.Buffer
	buff := bufio.NewWriter(&b)
	err := pdf.Output(buff)
	if err != nil {
		return nil, err
	}
	err = buff.Flush()
	return b.Bytes(), err
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
