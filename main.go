package main

import (
	"fmt"
	"time"

	cpupipeline "github.com/AndrewRazvanM/System-Monitor/systemmonitor/cpu_pipeline"
	layoutmanager "github.com/AndrewRazvanM/System-Monitor/systemmonitor/layout_manager"
	renderingsystem "github.com/AndrewRazvanM/System-Monitor/systemmonitor/rendering_system"
	ui "github.com/AndrewRazvanM/System-Monitor/systemmonitor/ui"
)

func main() {
	Buffer := renderingsystem.ScreenBuffer{}
	Buffer.Init(true, 4)

	rawReadings := cpupipeline.CPUReading{}
	if err := rawReadings.GetReady(); err != nil {
		fmt.Println("Ran into an error:\n", err)
	}

	widget := &ui.Widget{}
	widget.SetAttributes("Random", ui.Standard)

	cpuWidget := &ui.Widget{}
	cpuWidget.SetAttributes("CPU", ui.Blue)

	topRow := &layoutmanager.Split{
		Direction: layoutmanager.Row,
		Children: []layoutmanager.Entry{
			{Node: widget, Weight: 1.0},
		},
	}
	root := &layoutmanager.Split{
		Direction: layoutmanager.Column,
		Children: []layoutmanager.Entry{
			{Node: topRow},
			{Node: cpuWidget},
		},
	}
	root.UpdatePos(layoutmanager.Geometry{X: 0, Y: 0, W: Buffer.Width, H: Buffer.Height}, true)

	Buffer.ProcessCmds(widget.Draw())
	Buffer.ProcessCmds(cpuWidget.Draw())

	for {
		if lErr := rawReadings.GetCPULoad(); lErr != nil {
			fmt.Print(lErr, "\n")
		}
		if tErr := rawReadings.GetTemp(); tErr != nil {
			fmt.Println("Error getting temp readings:\n", tErr)
		}

		renderCommands := rawReadings.Compose(cpuWidget.Geometry)
		Buffer.ProcessCmds(renderCommands)

		Buffer.Render()
		Buffer.Screen.Show()
		time.Sleep(time.Second)
	}
}