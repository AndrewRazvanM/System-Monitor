package systemmonitor

import (
	"fmt"
)
const (
	cpuPadding int = 3 
	//this is the degrees symbol in ASCII - decimal
	degreesSymbol rune = '°'
)
// FormatReadins  takes the raw readings and formats them for the rendered.
// The formated values are sent to the Widget that's passed into it. 
func (cr *CPUReading) FormatReadings(cpuWidget *Widget) {
	currentIndex := 0
	for _, coreInfo := range cr.RawReadings {
		//converts to degress from milidegress
		tempDegrees := coreInfo.Temp / 1000
		//generates the actual line
		content := []rune(fmt.Sprintf("CPU %d %*d%cC", coreInfo.CoreID, cpuPadding, tempDegrees, degreesSymbol))
		for i, rune := range content {
			if i <= cpuWidget.Current.Width { 
			cpuWidget.Current.Cells[currentIndex + i].Rune = rune
			}
		}
		currentIndex += cpuWidget.Current.Width
	}
}