package systemmonitor

import (
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
	for i, rune := range w.Title {
		screen.SetContent(xTitle + i, y0, rune, nil, w.Style.Bold(true) )
	}

	//render the static elements on the dashboard
	header := "CPU     Temp°C Load"
	screen.PutStr(x0 + 1, y0 + 2, header)
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

// UpdateLine adds 1 line at a time to the widget.Current.Cells buffer. It will truncate the line if it's too long.
// It also assumes the caller overwrites the entire line.
func (w *Widget) UpdateLine(content []rune, line int,  ranges []StyleRange) {
	width := w.W - 2
	cellIndex := line * width
	cLen := len(content)
	rLen := len(ranges)
	maContentLen := len(w.Current.Cells)
	

	// if no style ranges are provided, build the content with default style.
	if rLen == 0 {
		for pos := 0; pos < width && pos < cLen; pos++ {
			cell := &w.Current.Cells[cellIndex + pos]
			cell.Rune = content[pos]
			cell.Style = Standard
		} 
		return
	}

	rIndex := 0
	r := ranges[rIndex]

	for pos := 0; pos < width && pos < cLen; pos++ {
		if cellIndex + pos >= maContentLen {
			break
		}
		// advance range if needed
		if rIndex < rLen - 1 && pos >= r.End {
			rIndex++
			r = ranges[rIndex]
		}

		cell := &w.Current.Cells[cellIndex + pos]
		cell.Rune = content[pos]

		// apply style if inside range
		if rIndex < rLen && pos >= r.Start && pos < r.End {
			cell.Style = r.Style
		} else {
			//tcell default style
			cell.Style = Standard
		}
	}
}

//Updates a single cell, based on it's index.
//Does not write outside the widget buffer.
func (w *Widget) UpdateCell(index int, char rune, style uint8) {

	if index < 0 || index >= len(w.Current.Cells) {
		return
	}

	w.Current.Cells[index].Rune = char
	w.Current.Cells[index].Style = style

}

func (w *Widget) Initalize(x, y, W, H int, title string) {
	w.X, w.Y, w.W, w.H = x, y, W, H
	w.IsVisible = true
	w.Style = tcell.StyleDefault
	w.Title = title
	w.Current.Cells = make([]Cell, (W - 2) * (H - 2))
	w.Previous.Cells = make([]Cell, (W - 2) * (H - 2))
}
