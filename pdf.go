package main

import (
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/signintech/gopdf"
)

type Color struct {
	R uint8
	G uint8
	B uint8
}

const (
	quantityColumnOffset = 360
	rateColumnOffset     = 405
	amountColumnOffset   = 480
)

const (
	subtotalLabel = "Subtotal"
	discountLabel = "Discount"
	taxLabel      = "Tax"
	totalLabel    = "Total"
	billToLabel   = "To"
	quantityLabel = "Days"
)

var text = Color{
	R: 25,
	G: 25,
	B: 25,
}
var border = Color{
	R: 225,
	G: 225,
	B: 225,
}
var heading = Color{
	R: 75,
	G: 75,
	B: 75,
}

func formatFloat(n float64) string {
	if n == float64(int64(n)) {
		return fmt.Sprintf("%.0f", n)
	}
	return fmt.Sprintf("%.1f", n)
}

func writeLogo(pdf *gopdf.GoPdf, logo string, from string) {
	if logo != "" {
		width, height := getImageDimension(logo)
		scaledWidth := 100.0
		scaledHeight := float64(height) * scaledWidth / float64(width)
		_ = pdf.Image(logo, pdf.GetX(), pdf.GetY(), &gopdf.Rect{W: scaledWidth, H: scaledHeight})
		pdf.Br(scaledHeight + 24)
	}
	_ = pdf.SetFont(fontNameBold, "", 14)
	pdf.SetTextColor(text.R, text.G, text.B)

	fromText := strings.Split(from, "\n")
	for i, line := range fromText {
		if i > 0 {
			_ = pdf.SetFont(fontName, "", 12)
		}
		_ = pdf.Cell(nil, line)
		pdf.Br(16)
	}
	pdf.Br(36)
	pdf.SetStrokeColor(border.R, border.G, border.B)
	pdf.Line(pdf.GetX(), pdf.GetY(), 100, pdf.GetY())
	pdf.Br(36)
}

func writeTitle(pdf *gopdf.GoPdf, title, id, date string) {
	_ = pdf.SetFont(fontNameBold, "", 24)
	pdf.SetTextColor(text.R, text.G, text.B)
	_ = pdf.Cell(nil, title)
	pdf.Br(36)
	_ = pdf.SetFont(fontName, "", 12)
	pdf.SetTextColor(heading.R, heading.G, heading.B)
	_ = pdf.Cell(nil, "#")
	_ = pdf.Cell(nil, id)
	pdf.SetTextColor(heading.R, heading.G, heading.B)
	_ = pdf.Cell(nil, "  ·  ")
	pdf.SetTextColor(heading.R, heading.G, heading.B)
	_ = pdf.Cell(nil, date)
	pdf.Br(48)
}

func writeDueDate(pdf *gopdf.GoPdf, due string) {
	_ = pdf.SetFont(fontName, "", 10)
	pdf.SetTextColor(heading.R, heading.G, heading.B)
	pdf.SetX(rateColumnOffset + 5)
	_ = pdf.CellWithOption(&gopdf.Rect{W: 45, H: 12}, "Due", gopdf.CellOption{Align: gopdf.Right | gopdf.Middle})
	pdf.SetTextColor(heading.R, heading.G, heading.B)
	_ = pdf.SetFontSize(10)
	pdf.SetX(amountColumnOffset - 15)
	_ = pdf.Cell(nil, due)
	pdf.Br(12)
}

func writeBillTo(pdf *gopdf.GoPdf, to string) {
	pdf.SetTextColor(heading.R, heading.G, heading.B)
	_ = pdf.SetFont(fontName, "", 9)
	_ = pdf.Cell(nil, strings.ToUpper(billToLabel))
	pdf.Br(18)
	pdf.SetTextColor(text.R, text.G, text.B)
	_ = pdf.SetFont(fontName, "", 12)
	texts := strings.Split(to, "\n")
	for _, text := range texts {
		_ = pdf.Cell(nil, text)
		pdf.Br(15)
	}
	pdf.Br(64)
}

func writeHeaderRow(pdf *gopdf.GoPdf) {
	_ = pdf.SetFont(fontName, "", 9)
	pdf.SetTextColor(heading.R, heading.G, heading.B)
	_ = pdf.Cell(nil, "ITEM")
	pdf.SetX(quantityColumnOffset)
	_ = pdf.Cell(nil, strings.ToUpper(quantityLabel))
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, "RATE")
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, "AMOUNT")
	pdf.Br(24)
}

func writeNotes(pdf *gopdf.GoPdf, notes string, header string) {
	pdf.SetY(650)

	_ = pdf.SetFont(fontName, "", 10)
	pdf.SetTextColor(heading.R, heading.G, heading.B)
	_ = pdf.Cell(nil, header)
	pdf.Br(18)
	_ = pdf.SetFont(fontName, "", 10)
	pdf.SetTextColor(text.R, text.G, text.B)

	formattedNotes := strings.ReplaceAll(notes, `\n`, "\n")
	notesLines := strings.Split(formattedNotes, "\n")

	for i := 0; i < len(notesLines); i++ {
		_ = pdf.Cell(nil, notesLines[i])
		pdf.Br(14)
	}

	pdf.Br(48)
}
func writeFooter(pdf *gopdf.GoPdf, id string) {
	pdf.SetY(800)

	_ = pdf.SetFont(fontName, "", 10)
	pdf.SetTextColor(150, 150, 150)
	_ = pdf.Cell(nil, id)
	pdf.SetStrokeColor(border.R, border.G, border.B)
	pdf.Line(pdf.GetX()+10, pdf.GetY()+6, 550, pdf.GetY()+6)
	pdf.Br(48)
}

func writeRow(pdf *gopdf.GoPdf, item string, quantity float64, rate float64) {
	_ = pdf.SetFont(fontName, "", 11)
	pdf.SetTextColor(text.R, text.G, text.B)

	total := float64(quantity) * rate

	_ = pdf.Cell(nil, item)
	pdf.SetX(quantityColumnOffset)
	_ = pdf.Cell(nil, formatFloat(quantity))
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, currencySymbols[file.Currency]+humanize.CommafWithDigits(rate, 2))
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, currencySymbols[file.Currency]+humanize.CommafWithDigits(total, 2))
	pdf.Br(24)
}

func writeTotals(pdf *gopdf.GoPdf, subtotal float64, tax float64, discount float64) {
	pdf.SetY(650)

	writeTotal(pdf, subtotalLabel, subtotal)
	if tax > 0 {
		writeTotal(pdf, taxLabel, tax)
	}
	if discount > 0 {
		writeTotal(pdf, discountLabel, discount)
	}
	writeTotal(pdf, totalLabel, subtotal+tax-discount)
}

func writeTotal(pdf *gopdf.GoPdf, label string, total float64) {
	_ = pdf.SetFont(fontName, "", 10)
	pdf.SetTextColor(heading.R, heading.G, heading.B)
	pdf.SetX(rateColumnOffset + 5)
	_ = pdf.CellWithOption(&gopdf.Rect{W: 45, H: 14}, label, gopdf.CellOption{Align: gopdf.Right | gopdf.Middle})
	pdf.SetTextColor(text.R, text.G, text.B)
	_ = pdf.SetFontSize(12)
	pdf.SetX(amountColumnOffset - 15)
	if label == totalLabel {
		_ = pdf.SetFont(fontNameBold, "", 11.5)
	}
	_ = pdf.Cell(nil, currencySymbols[file.Currency]+humanize.CommafWithDigits(total, 2))
	pdf.Br(24)
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
