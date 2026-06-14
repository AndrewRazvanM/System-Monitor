package systemmonitor

import "fmt"

const (
	cpuPadding = 3
	cpuDegrees = "°C"
)

func (cr *CPUReading) FormatReadings() {
	for core, rawTemp := range cr.RawTemp {
		tempDegrees := rawTemp / 1000
		cr.FormattedTemp[core] = fmt.Sprintf("%*d %s", cpuPadding, tempDegrees, cpuDegrees)
	}
}