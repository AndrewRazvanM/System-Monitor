package systemmonitor

import "fmt"

var barLoad = []rune("▏▎▍▌▋▊▉█")

const (
	CpuPadding int = 3 
	//this is the degrees symbol in ASCII - decimal
	degreesSymbol rune = '°'
	//it's where in the formated string the temp is listed. It's the len of
	//"Core" + CpuPadding + " "
	TempIndex int = len("CPU") + CpuPadding + 2
)
// FormatReadins  takes the raw readings and formats them for the rendered.
// The formated values are sent to the Widget that's passed into it. 
func (cr *CPUReading) UpdateBuffer(cpuWidget *Widget) {
	
	line := 0
	maxLine := cpuWidget.H - line - 2

	//This is the aggregate load style
	aLoadStyle := Standard
			switch {
			case cr.TotLoad.Load < 50:
				aLoadStyle = Green
			case cr.TotLoad.Load< 90:
				aLoadStyle = Yellow
			case cr.TotLoad.Load >= 90:
				aLoadStyle = Red
			}

	//This is the aggregated temp syle
	aTempStyle := Standard
			switch {
			case cr.TotLoad.Load < 50:
				aTempStyle = Green
			case cr.TotLoad.Load< 90:
				aTempStyle = Yellow
			case cr.TotLoad.Load >= 90:
				aTempStyle = Red
			}

	//This is the aggregated temp syle 
	n := cpuWidget.writeText(0 , maxLine - 1, "CPU", Standard)
	n += cpuWidget.writeText(n, maxLine - 1, fmt.Sprintf("%*d", CpuPadding, cr.TotLoad.Temp), aTempStyle) 
	n += cpuWidget.writeText(n, maxLine - 1, "°C ", Standard)
	n += cpuWidget.writeText(n, maxLine - 1, fmt.Sprintf("%*.1f%c", CpuPadding + 2, cr.TotLoad.Load, '%'), aLoadStyle)
	aggregateBarWidth := cpuWidget.W - 2 - n - 1
	n += cpuWidget.writeBar(n, maxLine - 1, aggregateBarWidth, cr.TotLoad.Load, aLoadStyle)

	//assuming all cores have the same number of thread -> get the max number of threads for core 0. Use that to calculate
	//the maximum available width for each thread.
	numOfMaxThreads := len(cr.RawReadings[0].Threads)
	// Each line is split in multiple parts; num of parts is determined by numOfMaxThreads. For example, if we're assuming 2 threads, the first half info for 1 thread, 2nd half info for another thread
	// Precomputing this value, so I don't calc it each time in the 2nd for loop
	availWidth := int(float64(cpuWidget.W - 2) / float64(numOfMaxThreads))

	for _, coreInfo := range cr.RawReadings {

		//converts to degress from milidegress
		tempDegrees := coreInfo.Temp / 1000
		tempStyle := Standard

		switch {
			case tempDegrees < 70:
				tempStyle = Green //green
			case tempDegrees < 85:
				tempStyle = Yellow //yellow
			case tempDegrees >= 85:
				tempStyle = Red //red
			}

		for threadNum, thread := range coreInfo.Threads {

			loadStyle := Standard
			switch {
			case thread.Load < 50:
				loadStyle = Green
			case thread.Load < 90:
				loadStyle = Yellow
			case thread.Load >= 90:
				loadStyle = Red
			}
			//tracks the position X in the line
			lineX := availWidth * threadNum
			x := cpuWidget.writeText(lineX , line, "CPU", Standard)
			x += cpuWidget.writeText(lineX + x, line, fmt.Sprintf("%*d", CpuPadding, thread.CPU), Standard)
			x += cpuWidget.writeText(lineX + x, line, fmt.Sprintf("%*d", CpuPadding, tempDegrees), tempStyle)
			x += cpuWidget.writeText(lineX + x, line, "°C", Standard)
			x += cpuWidget.writeText(lineX + x, line, fmt.Sprintf("%*.1f%c", CpuPadding + 2, thread.Load, '%'), loadStyle)
			barWidth := availWidth - x - 1
			x += cpuWidget.writeBar(lineX + x, line, barWidth, thread.Load, loadStyle)
			if threadNum >= numOfMaxThreads - 1 {
				continue
			}
			cpuWidget.UpdateCell(lineX + x, line, VLine, 4)
		}
		line++
		if line >= maxLine - 1 {
				break
		}
	}
}
