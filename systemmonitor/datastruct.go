package systemmonitor

import "os"

//store the temp data, caches the file handlers and stores the formatted strings (used by rendered)
type CPUReading struct {
		//store degress for eac core
		RawTemp map[string]int32
		//caches file descriptors for each core
		FilesDescriptor map[string]*os.File
		//stores the formatted strings that will be used by the renderer
		FormattedTemp map[string]string
	}
