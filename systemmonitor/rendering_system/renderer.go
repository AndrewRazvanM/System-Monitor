package systemmonitor

import (
	"fmt"

	tcell "github.com/gdamore/tcell/v2"
)

const (
    HLine = '─'
    VLine = '│'
    TLCorner = '┌'
    TRCorner = '┐'
    BLCorner = '└'
    BRCorner = '┘'
)

// DrawWidget draws the static elements of the Widget in it's style.
// Each widget stores it's own style.
func (w Widget) DrawWidget(screen tcell.Screen) {
	if !w.IsVisible {
		return
	}
	
	x0 := w.X
	y0 := w.Y
	x1 := w.X + w.W - 1
	y1 := w.Y + w.H - 1
	xTitle := x0 + 1

	style := w.Style

	// corners
	screen.SetContent(x0, y0, TLCorner, nil, style)
	screen.SetContent(x1, y0, TRCorner, nil, style)
	screen.SetContent(x0, y1, BLCorner, nil, style)
	screen.SetContent(x1, y1, BRCorner, nil, style)

	// top & bottom edges
	for x := x0 + 1; x < x1; x++ {
		screen.SetContent(x, y0, HLine, nil, style)
		screen.SetContent(x, y1, HLine, nil, style)
	}

	// left & right edges
	for y := y0 + 1; y < y1; y++ {
		screen.SetContent(x0, y, VLine, nil, style)
		screen.SetContent(x1, y, VLine, nil, style)
	}

	//render title
	title := fmt.Sprintf("%c%s%c", TRCorner, w.Title, TLCorner)
	for i, rune := range title {
		screen.SetContent(xTitle + i, y0, rune, nil, w.Style.Bold(true) )
	}

}

// Render adds the content into tcell's buffer.
func (w *Widget) Render(screen tcell.Screen) {
	for y := 0; y < w.H - 2; y++ {
		for x := 0; x < w.W - 2; x++ {
			index := y * (w.W - 2) + x
			//render only if content is different
			if w.Previous.Cells[index] == w.Current.Cells[index] {
				continue
			}
			cell := w.Current.Cells[index]
			screen.SetContent(x + w.X + 1, y + w.Y + 1, cell.Rune, nil, StyleList[cell.Style])

			w.Previous.Cells[index] = cell
		}	
	}	
}

//Updates a single cell, based on it's index.
//Does not write outside the widget buffer.
func (w *Widget) UpdateCell(x, y int, char rune, style uint8) {
	index := y * (w.W - 2) + x
	if index < 0 || index >= len(w.Current.Cells) {
		return
	}

	w.Current.Cells[index].Rune = char
	w.Current.Cells[index].Style = style

}

func (w *Widget) Initalize(x, y, W, H int, title string, style uint8) {
	w.X, w.Y, w.W, w.H = x, y, W, H
	w.IsVisible = true
	w.Style = StyleList[style]
	w.Title = title
	w.Current.Cells = make([]Cell, (W - 2) * (H - 2))
	w.Previous.Cells = make([]Cell, (W - 2) * (H - 2))
}

// writeText takes a list of runes and places it into the correct place in the Cells buffer.
// It returns the number of cell that were written into. 
// It also checks if the []rune is bigger than the buffer, if it is, it returns a 0.
func (w *Widget) writeText(x, y int, text string, style uint8) (written int) {
	cellIndex := y*(w.W-2) + x
	max := len(w.Current.Cells)

	for _, r := range text {
		if cellIndex + written >= max {
			break
		}

		w.Current.Cells[cellIndex + written].Rune = r
		w.Current.Cells[cellIndex + written].Style = style
		written++
	}

	return written
}
// writeBar writes a bar in the widget buffer, with the given character (barChar).
// It also checks if the width will write outside the buffer, if it will, it returns a 0.
func (w *Widget) writeBar(x, y int, width int, pct float64, style uint8) (cellsWritten int) {
	cellIndex := y * (w.W - 2) + x
	maxCellLen := len(w.Current.Cells)
	
	if cellIndex < 0 || cellIndex + width > maxCellLen {
		return 0
	}
	filled := int(pct * float64(width) / 100.0)

    for i := range width{
		cell := w.Current.Cells[cellIndex + i]
        c := '─'
        s := Standard

        if i < filled {
            c = '━'
			s = style
        }

		cell.Rune = c
		cell.Style = s
		w.Current.Cells[cellIndex + i] = cell

    }
	return width
}

func (w *Widget) fillRect(x, y, width, height int, char rune, style uint8,) {

}
