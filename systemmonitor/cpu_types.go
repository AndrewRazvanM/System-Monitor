package systemmonitor

import "os"

//ThreadInfo is used to track the cpu load per Thread
type ThreadInfo struct {
    CPU  int
    Load float32

	//stores previous raw readings
	//used to calculate the load
	prevBusy uint64
	prevTotal uint64
}

//CoreInfo is used to track the per core information.
//Each core can have multiple Threads
type CoreInfo struct { 
    CoreID int
    Temp int32
	//Threads is a map[ThreadID]ThreadInfo. Each ThreadInfo is saved by it's CPU ID (Thread ID)
    Threads []ThreadInfo
	//caches the open file
	File *os.File
}

type AggregateLoad struct {
	//Stores the total CPU load for the system
	Load float32
	//Stores the previous time the CPU spent doing work.
	//It's used to calculate the total Load delta
	prevBusyTime uint64
	//Stores the previous time the CPU spent doing work + time spent idle.
	//It's used to calculate the total Load delta
	prevTotTime uint64
}

// CPUReading store the CoreInfo and caches the file handlers.
type CPUReading struct {
		//Store the load for the whole system and the values used to calculate it.
		TotLoad AggregateLoad

		//Store degress for eac core
		RawReadings []CoreInfo
		//Stores the CPU topography. It maps physical cores to threads.
		//Each physical core is mapped to a list of threads.
		CPUTopology map[int][]int
		//Stores a reverse map of the CPU topology.
		//It's used by the GetCpuLoad() to efficiently map the thread readings to the physical cores
		rCPUTopology map[int]int
		//Stores the /proc/stat file.
		//This gives me the CPU Load per thread
		StatFile *os.File
	}
