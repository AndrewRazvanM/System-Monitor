package systemmonitor

import tcell "github.com/gdamore/tcell/v2"

const (
    HLine = '─'
    VLine = '│'
    TLCorner = '┌'
    TRCorner = '┐'
    BLCorner = '└'
    BRCorner = '┘'
)

// DrawWidget draws the borders for the Widget in it's style.
// Each widget stores it's own style.
func (w Widget) DrawWidget(screen tcell.Screen) {
	if !w.IsVisible {
		return
	}

	x0 := w.X
	y0 := w.Y
	x1 := w.X + w.W - 1
	y1 := w.Y + w.H - 1

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
}

func (w Widget) Render(screen tcell.Screen) {
	for y := 0; y < w.H - 2; y++ {
		for x := 0; x < w.W - 2; x++ {
			cell := w.Current.Cells[y * (w.W - 2) + x]
			screen.SetContent(x + w.X + 1, y + w.Y + 1, cell.Rune, nil, tcell.StyleDefault)
		}
	}
}

func (w *Widget) Setline(content []rune, line int,  ranges []StyleRange) {
	width := w.W - 2
	cellIndex := line * width
	cLen := len(content)
	rLen := len(ranges) 
	

	// if no style ranges are provided, build the content with default style.
	if rLen == 0 {
		for pos := 0; pos < width && pos < cLen; pos++ {
			cell := &w.Current.Cells[cellIndex + pos]
			cell.Rune = content[pos]
			cell.Style = tcell.StyleDefault
		}
		return
	}

	rIndex := 0
	r := ranges[rIndex]

	for pos := 0; pos < width && pos < cLen; pos++ {

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
			cell.Style = tcell.StyleDefault
		}
	}
}

func (w *Widget) Initalize(x, y, W, H int, title string) {
	w.X, w.Y, w.W, w.H = x, y, W, H
	w.IsVisible = true
	w.Style = tcell.StyleDefault
	w.Title = title
	w.Current.Cells = make([]Cell, (W - 2) * (H - 2))
}