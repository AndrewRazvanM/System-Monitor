package systemmonitor

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)
var (
    CoolTempStyle = tcell.StyleDefault.
        Foreground(tcell.ColorGreen)

    WarmTempStyle = tcell.StyleDefault.
        Foreground(tcell.ColorYellow)

    HotTempStyle = tcell.StyleDefault.
        Foreground(tcell.ColorRed)
)
const (
	cpuPadding int = 3 
	//this is the degrees symbol in ASCII - decimal
	degreesSymbol rune = '°'
	headerPadding int = 2
	tempIndex int = 7
)
// FormatReadins  takes the raw readings and formats them for the rendered.
// The formated values are sent to the Widget that's passed into it. 
func (cr *CPUReading) FormatReadings(cpuWidget *Widget) {
	line := 0 + headerPadding
	maxLine := cpuWidget.H - headerPadding - 1
	// Tells the widget what style secific parts of the line should have.
	// This should be ordered: First to Last
	arrayStyles := make([]StyleRange, 3)

	for _, coreInfo := range cr.RawReadings {
		//converts to degress from milidegress
		tempDegrees := coreInfo.Temp / 1000
		switch {
		case tempDegrees < 70:
			arrayStyles[0] = StyleRange{
				Start: tempIndex,
				End: tempIndex + cpuPadding,
				Style: CoolTempStyle,
			}
		case tempDegrees < 85:
			arrayStyles[0] = StyleRange{
				Start: tempIndex,
				End: tempIndex + cpuPadding,
				Style: WarmTempStyle,
		}
		default:
			arrayStyles[0] = StyleRange{
				Start: tempIndex,
				End: tempIndex + cpuPadding,
				Style: WarmTempStyle,
		}
		}
		//generates the actual line
		content := []rune(fmt.Sprintf("CPU%*d %*d%cC", cpuPadding, coreInfo.CoreID, cpuPadding, tempDegrees, degreesSymbol))
		cpuWidget.Setline(content, line, arrayStyles)

		line++
		if line > maxLine {
			break
		}
	}
}