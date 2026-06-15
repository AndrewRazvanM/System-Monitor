package systemmonitor

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

type CoreTemp struct {
	Name string
	Temp int8
}
// CPUReading store the temp data and caches the file handlers.
type CPUReading struct {
		//store degress for eac core
		RawTemps map[string]int32
		//caches file descriptors for each core
		FilesDescriptor map[string]*os.File
	}


// Widget stores the position and dimensions for 
// each dashboard used by the renderer
type Widget struct {
	//window positions 
	X, Y int
	//window size
	W, H int
	//says if the widget should show or not
	IsVisible bool
	//stores the style used by tcell
	Style tcell.Style
	//widget title
	Title string
	//stores the formatted strings that will be used by the renderer
	FormattedTemp map[string]string
}