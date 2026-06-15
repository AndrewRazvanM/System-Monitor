package main

import (
	"fmt"

	"github.com/AndrewRazvanM/System-Monitor/systemmonitor"
	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, sErr := tcell.NewScreen()
	if sErr != nil {
		fmt.Println("error generating window")
	}

	initErr := screen.Init()
	
	if initErr != nil {
		fmt.Println("error initializing screen ", initErr)
	}
	rawReadings := systemmonitor.CPUReading{}
	widget := systemmonitor.Widget{}
	widget.Initalize(5, 5, 51, 12, "CPU Dashboard")

	err := rawReadings.GetReady()
	if err != nil {
		fmt.Println("Ran into an error:\n", err)
	}
	tErr := rawReadings.GetTemp()
	if tErr != nil {
		fmt.Println("Error getting temp readings:\n", tErr)
	}
	rawReadings.FormatReadings(&widget)
	widget.DrawWidget(screen)
	widget.Render(screen)
	screen.Show()
	}