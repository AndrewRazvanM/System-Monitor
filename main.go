package main

import (
	"fmt"
	"time"

	cpupipeline "github.com/AndrewRazvanM/System-Monitor/systemmonitor/cpu_pipeline"
	renderingsystem "github.com/AndrewRazvanM/System-Monitor/systemmonitor/rendering_system"
	ui "github.com/AndrewRazvanM/System-Monitor/systemmonitor/ui"
)

func main() {
	Buffer := renderingsystem.ScreenBuffer{}
	Buffer.Init(true, 4)
	rawReadings := cpupipeline.CPUReading{}
	widget := ui.Widget{}
	widget.Initalize(10, 3, 100, 12, "CPU Dashboard", 4)
	widgetCmds := widget.Draw()
	Buffer.ProcessCmds(widgetCmds)

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
	
	renderCommands := rawReadings.Compose(widget.Geometry)

	Buffer.ProcessCmds(renderCommands)
	Buffer.Render()
	time.Sleep(1000000000)
	}
}