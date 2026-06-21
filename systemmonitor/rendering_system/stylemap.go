package renderingsystem

import (
	"github.com/AndrewRazvanM/System-Monitor/systemmonitor/ui"
	"github.com/gdamore/tcell/v2"
)

var StyleList = []tcell.Style{
	ui.Green: tcell.StyleDefault.Foreground(tcell.ColorGreen),  // 0 - Green
	ui.Yellow: tcell.StyleDefault.Foreground(tcell.ColorYellow), // 1 - Yellow
	ui.Red: tcell.StyleDefault.Foreground(tcell.ColorRed),    // 2 - Red
	ui.Standard: tcell.StyleDefault,								  // 3 - standard
	ui.Blue: tcell.StyleDefault.Foreground(tcell.ColorBlue),	  // 4 - Blue
}

var StyleAttr = []tcell.AttrMask {
	ui.NoAttr: tcell.AttrNone,
	ui.Bold:	tcell.AttrBold,
	ui.Dim:	tcell.AttrDim,
}