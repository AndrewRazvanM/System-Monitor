package systemmonitor

import (
	"fmt"
)

var barLoad = []rune("▏▎▍▌▋▊▉█")

const (
	cpuPadding int = 3 
	//this is the degrees symbol in ASCII - decimal
	degreesSymbol rune = '°'
	headerPadding int = 2
	//it's where in the formated string the temp is listed. It's the len of
	//"Core" + cpuPadding + " "
	tempIndex int = len("CPU") + cpuPadding + 2
)
// FormatReadins  takes the raw readings and formats them for the rendered.
// The formated values are sent to the Widget that's passed into it. 
func (cr *CPUReading) FormatReadings(cpuWidget *Widget) {
	
	line := 0 + headerPadding
	maxLine := cpuWidget.H - line - 1

	//assuming all cores have the same number of thread -> get the max number of threads for core 0. Use that to calculate
	//the maximum available width for each thread.
	numOfMaxThreads := len(cr.RawReadings[0].Threads)
	// Each line is split in multiple parts; num of parts is determined by numOfMaxThreads. For example, if we're assuming 2 threads, the first half info for 1 thread, 2nd half info for another thread
	// Precomputing this value, so I don't calc it each time in the 2nd for loop
	availWidth := float32(cpuWidget.W - 2) / float32(numOfMaxThreads)

	for _, coreInfo := range cr.RawReadings {
		// Tells the widget what style secific parts of the line should have.
		// This should be ordered: First to Last
		// there will be a max of 4 array elements per dashboard line
		arrayStyles := make([]StyleRange, numOfMaxThreads * 2)

		//converts to degress from milidegress
		tempDegrees := coreInfo.Temp / 1000
	
		//index for the style array. Keeps track of how many are added
		ind := 0
		//stores all the output for each thread
		threadLines := make([]rune, 0 , cpuWidget.W - 2)
		//generates the other 2 halfs
		for threadNum, thread := range coreInfo.Threads {
			message := []rune(fmt.Sprintf("CPU%*d %*d%cC %*.1f%c ", cpuPadding, thread.CPU, cpuPadding, tempDegrees, degreesSymbol, cpuPadding + 1, thread.Load, '%'))
			
			//tracks where the load is in the string
			//the message is created with consisten padding. The load should always be in the last 5 char of the message
			indLoad := len(threadLines) + len(message) - 6
			//generates the thread load bar
			barWidth := int(thread.Load / 100 * availWidth)

			//accounting for the borders + the len of message. I need to leave space for the next threads
			for i := range (cpuWidget.W - 2) / 2 - len(message) {
				if i <= barWidth - 1{
					message = append(message, '━')
				} else {
					message = append(message, '─')
				}
			}
			//add divider between threads
			message = append(message, '|')
			threadLines = append(threadLines, message...)

			//assign the style for the temp and load:
			//for temp
			switch {
			case tempDegrees < 70:
				arrayStyles[ind] = StyleRange{
					Start: tempIndex + threadNum * ((cpuWidget.W - 2) / 2),
					End: tempIndex + threadNum * ((cpuWidget.W - 2) / 2) + cpuPadding,
					Style: Green, //green
				}
			case tempDegrees < 85:
				arrayStyles[ind] = StyleRange{
					Start: tempIndex + threadNum * ((cpuWidget.W - 2) / 2),
					End: tempIndex + threadNum * ((cpuWidget.W - 2) / 2) + cpuPadding,
					Style: Yellow, //yellow
				}
			case tempDegrees >= 85:
				arrayStyles[ind] = StyleRange{
					Start: tempIndex + threadNum * ((cpuWidget.W - 2) / 2),
					End: tempIndex + threadNum * ((cpuWidget.W - 2) / 2) + cpuPadding,
					Style: Red, //red
				}
			}

			ind++

			switch {
			case thread.Load < 50:
				arrayStyles[ind] = StyleRange{
					Start: indLoad,
					//covers the actual cpu load and the percentage sign
					End: indLoad + cpuPadding + barWidth,
					Style: Green, //green
				}
			case thread.Load < 90:
				arrayStyles[ind] = StyleRange{
					Start: indLoad,
					End: indLoad + cpuPadding + barWidth,
					Style: Yellow, //yellow
				}
			default:
				arrayStyles[ind] = StyleRange{
					Start: indLoad,
					End: indLoad + cpuPadding + barWidth,
					Style: Red, //red
				}
			}
			ind++
		}
		cpuWidget.UpdateLine(threadLines, line, arrayStyles)
		line++
		if line >= maxLine {
				break
		}
	}
}