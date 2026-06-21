package renderingsystem

import (
	ui "github.com/AndrewRazvanM/System-Monitor/systemmonitor/ui"
	"github.com/gdamore/tcell/v2"
)

//ScreenBuffer is the entire screen
type ScreenBuffer struct {	
	//Holds current buffer
	Current Snapshot
	//Holds previous buffer
	Previous Snapshot
	//Width of the screen
	Width int
	//Height of the screen
	Height int
	//Holds the screen that's associated with this buffer
	Screen tcell.Screen
	//Might have multiple screens in the future
	IsVisible bool
}

// Cell stores the actual character (element) that will be rendered.
// Also contains the position it should be rendered at and it's style
type Cell struct {
    Rune rune
    Style ui.Color
	Attribute ui.CellAttr
}

// Snapshot stores the elements that will be rendered.
// Each Snapshot is a 2D Grid mapped on a slice (1D).
type Snapshot struct {
	// Slice of 2D positions (x,y). Basically mapping 2D to 1D:
	// Cells[Index] := Y * W + X 
    Cells []Cell
}
