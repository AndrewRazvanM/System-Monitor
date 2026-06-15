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

// DrawWidget draws the borders for the Widget in a specific style.
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

func GetReady() {
	
}