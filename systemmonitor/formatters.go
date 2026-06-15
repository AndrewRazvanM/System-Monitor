package systemmonitor

import (
	"fmt"
)
const (
	cpuPadding = 3
	cpuDegrees = "°C"
)

func (cr *CPUReading) FormatReadings(widget Widget) {
	for core, rawTemp := range cr.RawTemps {
		tempDegrees := rawTemp / 1000
		widget.FormattedTemp[core] = fmt.Sprintf("%*d %s", cpuPadding, tempDegrees, cpuDegrees)
	}
}