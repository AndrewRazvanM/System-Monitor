package renderingsystem

import "github.com/gdamore/tcell/v2"

var StyleList = []tcell.Style{
	tcell.StyleDefault.Foreground(tcell.ColorGreen),  // 0 - Green
	tcell.StyleDefault.Foreground(tcell.ColorYellow), // 1 - Yellow
	tcell.StyleDefault.Foreground(tcell.ColorRed),    // 2 - Red
	tcell.StyleDefault,								  // 3 - standard
	tcell.StyleDefault.Foreground(tcell.ColorBlue),	  // 4 - Blue
}

