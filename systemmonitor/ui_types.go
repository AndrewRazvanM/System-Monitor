package systemmonitor

import (
	"github.com/gdamore/tcell/v2"
)

// Geometry stores the pos on screen (x, y) and size (W, H)
type Geometry struct {
	//window positions 
	X, Y int
	// Width and Height of Widget
	W, H int
}

// Widget stores the position and dimensions for 
// each dashboard used by the renderer
type Widget struct {
	//the Widget knows where it is, cause it knows where it isn't
	//only the layout manager updates this
	Geometry
	//IsVisible says if the widget should show or not
	IsVisible bool
	//Dirty says if the widget needs to be redrawn
	Dirty bool
	//Style stores the style used by the widget
	Style tcell.Style
	//Title contains the widget title
	Title string

	//Current stores the current buffer used by the renderer
	Current  Snapshot
	//Previous stores the previous buffer, used to check if the content changed
    Previous Snapshot

	//Formater takes the raw data and formats it into the correct position in the Current Snapshot
	Formater func(*Widget)
}

// StyleRange is used by the formater. It's used to tell the widget
// what style to set for specific parts of the line.
type StyleRange struct {
    Start int
    End   int
    Style uint8
}