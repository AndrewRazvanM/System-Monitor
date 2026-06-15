package systemmonitor

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

//ThreadInfo is used to track the cpu load per Thread
type ThreadInfo struct {
    CPU  int
    Load float32
}

//CoreInfo is used to track the per core information.
//Each core can have multiple Threads
type CoreInfo struct { 
    CoreID int
    Temp int32
    Threads []ThreadInfo
}

// CPUReading store the CoreInfo and caches the file handlers.
type CPUReading struct {
		//store degress for eac core
		RawReadings []CoreInfo
		//caches file descriptors for each core
		FilesDescriptor map[int]*os.File
		//stores the CPU topography. It maps physical cores to threads.
		//each physical core is mapped to a list of threads.
		CPUTopology map[int][]int
	}

// Cell stores the actual character (element) that will be rendered.
// Also contains the position it should be rendered at and it's style
type Cell struct {
    Rune rune
    Style tcell.Style
}

// Snapshot stores the elements that will be rendered.
type Snapshot struct {
	Width, Height int
    Cells []Cell
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
	//says if the widget needs to be redrawn
	Dirty bool
	//stores the style used by tcell
	Style tcell.Style
	//widget title
	Title string
	//stores snapshts to be used by the renderer and the content diff engine
	Current  Snapshot
    Previous Snapshot
}