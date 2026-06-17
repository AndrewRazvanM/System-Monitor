package systemmonitor

import "github.com/gdamore/tcell/v2"

var StyleList = []tcell.Style{
	tcell.StyleDefault.Foreground(tcell.ColorGreen),  // 0 - Green
	tcell.StyleDefault.Foreground(tcell.ColorYellow), // 1 - Yellow
	tcell.StyleDefault.Foreground(tcell.ColorRed),    // 2 - Red
	tcell.StyleDefault,								  // 3 - standard
}

const (
	Green = 0
	Yellow = 1
	Red = 2
	Standard = 3 // maps to tcell.StyleDefault
)