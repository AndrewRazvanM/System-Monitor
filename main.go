package main

import (
	"fmt"

	"github.com/AndrewRazvanM/System-Monitor/systemmonitor"
	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("Ran into error when initializing the screen:\n", err)
	}
	if err := screen.Init(); err != nil {
    fmt.Println(err)
    return
}

	style := tcell.StyleDefault
	tWidget := systemmonitor.Widget{
    X: 2,
    Y: 6,
    W: 51,
    H: 10,
    IsVisible: true,
    Style: style,
	}	
	tWidget.DrawWidget(screen)
	screen.Show()
	}