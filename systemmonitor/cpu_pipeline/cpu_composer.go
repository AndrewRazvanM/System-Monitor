package cpupipeline

import (
	"fmt"
	"unicode/utf8"

	ui "github.com/AndrewRazvanM/System-Monitor/systemmonitor/ui"
)


var barLoad = []rune("▏▎▍▌▋▊▉█")
const (
	eSkinnyBar rune = '─'
	fSkinnyBar rune = '━'
)

const (
	CpuPadding int = 3 
	//this is the degrees symbol in ASCII - decimal
	degreesSymbol rune = '°'
)


func (cr *CPUReading) Compose(area ui.Geometry) []ui.DrawCommand {
	if len(cr.RawReadings) == 0 {
		return nil
	}

	cmds := make([]ui.DrawCommand, 0, 256)

	x0 := area.X + 1
	y := area.Y + 1

	maxWidth := area.W - 2
	maxLine := area.Y + area.H - 2

	numThreads := len(cr.RawReadings[0].Threads)
	if numThreads == 0 {
		return cmds
	}

	maxBarWidth := (area.W - 2) / numThreads

	// --------------------------------------------------
	// Aggregate row (bottom of widget)
	// --------------------------------------------------

	aggY := area.Y + area.H - 2

	loadStyle := cpuLoadStyle(cr.TotLoad.Load)
	tempStyle := cpuTempStyle(cr.TotLoad.Temp / 1000)

	x := 0

	cmds = append(cmds,
		ui.DrawCommand{
			Type:  ui.CommandText,
			X:     x0 + x,
			Y:     aggY,
			Text:  "CPU ",
			Style: ui.Standard,
		},
	)
	x += 4

	tempStr := fmt.Sprintf("%*d°C ", CpuPadding, cr.TotLoad.Temp/1000)
	cmds = append(cmds, ui.DrawCommand{
		Type:  ui.CommandText,
		X:     x0 + x,
		Y:     aggY,
		Text:  tempStr,
		Style: tempStyle,
	})
	x += utf8.RuneCountInString(tempStr)

	loadStr := fmt.Sprintf("%*.1f%%", CpuPadding+2, cr.TotLoad.Load)
	cmds = append(cmds, ui.DrawCommand{
		Type:  ui.CommandText,
		X:     x0 + x,
		Y:     aggY,
		Text:  loadStr,
		Style: loadStyle,
	})
	x += utf8.RuneCountInString(loadStr)

	cmds = append(cmds, ui.DrawCommand{
		Type:  ui.CommandFill,
		X:     x0 + x,
		Y:     aggY,
		W:     maxWidth - x,
		H:     1,
		Char:  fSkinnyBar,
		Style: loadStyle,
	})

	// --------------------------------------------------
	// Core rows
	// --------------------------------------------------

	for _, core := range cr.RawReadings {
		if y >= maxLine {
			break
		}

		x = 0

		tempDeg := core.Temp / 1000
		coreTempStyle := cpuTempStyle(tempDeg)

		for j, thread := range core.Threads {

			loadStyle := cpuLoadStyle(thread.Load)

			cpuStr := fmt.Sprintf("%-*d", CpuPadding, thread.CPU)
			tempStr := fmt.Sprintf("%*d°C", CpuPadding, tempDeg)
			loadStr := fmt.Sprintf("%*.1f%%", CpuPadding+2, thread.Load)

			cmds = append(cmds,
				ui.DrawCommand{
					Type:  ui.CommandText,
					X:     x0 + x,
					Y:     y,
					Text:  cpuStr,
					Style: ui.Standard,
				},
			)
			x += utf8.RuneCountInString(cpuStr)

			cmds = append(cmds,
				ui.DrawCommand{
					Type:  ui.CommandText,
					X:     x0 + x,
					Y:     y,
					Text:  tempStr,
					Style: coreTempStyle,
				},
			)
			x += utf8.RuneCountInString(tempStr)

			cmds = append(cmds,
				ui.DrawCommand{
					Type:  ui.CommandText,
					X:     x0 + x,
					Y:     y,
					Text:  loadStr,
					Style: loadStyle,
				},
			)
			x += utf8.RuneCountInString(loadStr)

			barWidth := (maxBarWidth * (j + 1)) - x

			cmds = append(cmds, ui.DrawCommand{
				Type:  ui.CommandFill,
				X:     x0 + x,
				Y:     y,
				W:     barWidth,
				H:     1,
				Char:  fSkinnyBar,
				Style: loadStyle,
			})

			x += barWidth

			// separator
			if j < numThreads-1 && x < maxWidth {
				cmds = append(cmds, ui.DrawCommand{
					Type:  ui.CommandRune,
					X:     x0 + x,
					Y:     y,
					Char:  ui.VLine,
					Style: ui.Blue,
				})
				x++
			}
		}

		y++
	}

	return cmds
}

func cpuLoadStyle (cpuLoad float64) ui.Color {
	style := ui.Standard
	switch {
		case cpuLoad < 50:
				style = ui.Green
		case cpuLoad < 90:
				style = ui.Yellow
		case cpuLoad >= 90:
				style = ui.Red
		}
	return style
}

func cpuTempStyle (temp int32) ui.Color {
	style := ui.Standard
	switch {
		case temp < 72:
				style = ui.Green
		case temp < 85:
				style = ui.Yellow
		case temp >= 85:
				style = ui.Red
		}
	return style
}
