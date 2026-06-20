package main

import (
	"fmt"
	"time"

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
	widget.Initalize(0, 0, 50, 12, "CPU Dashboard", 4)

	err := rawReadings.GetReady()
	if err != nil {
		fmt.Println("Ran into an error:\n", err)
	}
	
	for {
	lErr := rawReadings.GetCPULoad()
	if lErr != nil {
		fmt.Print(lErr, "\n")
	}
	tErr := rawReadings.GetTemp()
	if tErr != nil {
		fmt.Println("Error getting temp readings:\n", tErr)
	}
	rawReadings.UpdateBuffer(&widget)

	widget.DrawWidget(screen)
	widget.Render(screen)
	screen.Show()
	time.Sleep(1000000000)
	}
}