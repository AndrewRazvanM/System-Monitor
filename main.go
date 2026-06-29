package main

import (
	"fmt"
	"time"

	cpupipeline "github.com/AndrewRazvanM/System-Monitor/systemmonitor/cpu_pipeline"
	inputhandler "github.com/AndrewRazvanM/System-Monitor/systemmonitor/input_handler"
	layoutmanager "github.com/AndrewRazvanM/System-Monitor/systemmonitor/layout_manager"
	renderingsystem "github.com/AndrewRazvanM/System-Monitor/systemmonitor/rendering_system"
	ui "github.com/AndrewRazvanM/System-Monitor/systemmonitor/ui"
)

func main() {
	buffer := renderingsystem.ScreenBuffer{}
	buffer.Init(true, ui.Standard)

	rawReadings := cpupipeline.CPUReading{}
	if err := rawReadings.GetReady(); err != nil {
		fmt.Println(err)
		return
	}

	// Widgets
	widget := &ui.Widget{}
	widget.SetAttributes("Random", ui.Standard)

	memWidget := &ui.Widget{}
	memWidget.SetAttributes("Memory", ui.Green)

	cpuWidget := &ui.Widget{}
	cpuWidget.SetAttributes("CPU", ui.Blue)

	// Layout
	topRow := &layoutmanager.Split{
		Direction: layoutmanager.Row,
		Children: []layoutmanager.Entry{
			{Node: widget},
			{Node: memWidget},
			{Node: cpuWidget},
		},
	}

	root := layoutmanager.Split{
		Direction: layoutmanager.Column,
		Children: []layoutmanager.Entry{
			{Node: topRow},
		},
	}

	root.UpdatePos(layoutmanager.Geometry{
		X: 0,
		Y: 0,
		W: buffer.Width,
		H: buffer.Height,
	}, true)

	buffer.ProcessCmds(widget.Draw())
	buffer.ProcessCmds(memWidget.Draw())
	buffer.ProcessCmds(cpuWidget.Draw())

	// Test button
	cpuWidget.Buttons = append(cpuWidget.Buttons, ui.Button{
		Label: "Disable",
		Geometry: layoutmanager.Geometry{
			X: cpuWidget.X,
			Y: cpuWidget.Y,
			W: 11,
			H: 1,
		},
		Visible: true,
		Behavior: func() {
			ui.ToggleWidget(cpuWidget)
		},
	})

	// Input handler
	input := inputhandler.InputHandler{
		WidgetButtons: []*ui.Widget{
			widget,
			memWidget,
			cpuWidget,
		},
	}

	go input.HandleInput(0, 0, root, buffer)

	for {
		if err := rawReadings.GetCPULoad(); err != nil {
			fmt.Println(err)
		}

		if err := rawReadings.GetTemp(); err != nil {
			fmt.Println(err)
		}

		buffer.ProcessCmds(rawReadings.Compose(cpuWidget.Geometry, cpuWidget.IsVisible))

		buffer.Render()

		time.Sleep(time.Second)
	}
}