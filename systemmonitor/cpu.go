package systemmonitor

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Each time a new reading is needed and it's file is cached, this function has to be used.
// Resets the file to pos 0.
func resetOpenFile (file *os.File) error {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}

func (cr CPUReading) findFolder() (string, error) {
	const cpuTempFolderPath = "/sys/class/hwmon"
	//sensors ordered by priority. Ideally, coretemp (Intel) or k10temp (AMD) are available.
	sensors := [6]string{"coretemp", "k10temp", "x86_pkg_temp", "cpu_thermal", "tcpu", "acpitz"}

	folderList, err := os.ReadDir(cpuTempFolderPath)
	if err != nil {
		return "", err
	}

	for _, sensor := range sensors {

		for _, v := range folderList {
		sensorFolderName := v.Name()
		sensorNamePath := filepath.Join(cpuTempFolderPath, sensorFolderName, "name")
		data, err := os.ReadFile(sensorNamePath)
		if err != nil {
			continue
		}
		strData := strings.TrimSpace(string(data))
		//exit on first sensor hit. From Best -> Worst
		if sensor == strData {
			filePath := filepath.Join(cpuTempFolderPath, sensorFolderName)
			return filePath, nil
		}
		}
	}
	return "", errors.New("Could not find a sensor")
}

// Close goes over all open files and closes them.
// All errors are captured and returned as a list. If no errors, the list is initialized with the 0 value.
func (cr *CPUReading) CloseFiles() []error {
	err := []error{}
	for _, v := range cr.RawReadings {
		error := v.File.Close()
		if error != nil {
			errorMsg := fmt.Sprintln("File ", v.File.Name(), " could not be closed. Error: ", error)
			err = append(err, errors.New(errorMsg))
		}
	}

	return err
}

// mapTempFiles goes through the sensor folder, checks the temp label file and opens the temp input file.
// It creates a map with the label (key) and the open file (element).
// Each time new readings are needed, the stored files are looped over.
func (cr *CPUReading) mapTempFiles(rootPath string) error {
	fileList, rError := os.ReadDir(rootPath)
	if rError != nil {
		return rError
	}
	
	for _, v := range fileList {
		//get the label name
		fileName := v.Name()
		//add open files to CpuReading.FileDescriptors
		if strings.Contains(fileName, "temp") && strings.Contains(fileName, "label") {
			labelPath := filepath.Join(rootPath, fileName)
			data, dError := os.ReadFile(labelPath)
			if dError != nil {
				return dError
			}

			//convert the data from bytes of ASCII number to string
			strData := strings.TrimSpace(string(data))
			//find the indeX of the last space in the string -> create a slice from the string,
			//starting from index + 1 -> convert this to an int
			label, intErr := strconv.Atoi(strData[strings.LastIndex(strData, " ") + 1:])
			if intErr != nil {
				return intErr
			}

			//check for the input file -> it contains the raw temp 
			inputStr := strings.TrimSuffix(fileName, "label") + "input"
			inputPath := filepath.Join(rootPath, inputStr)
			FileDescriptor, iError := os.Open(inputPath)
			if iError != nil {
				return iError
			}
			//convert label from int64 to int. 
			cr.RawReadings[label].File = FileDescriptor
		}
	}
	return nil
}

func (cr *CPUReading) getCPUTopology() error {
	const rootPath = "/sys/devices/system/cpu"
	const topologyFilePath = "topology/core_id"
	rootFolderList, fErr := os.ReadDir(rootPath)
	if fErr != nil {
		return fErr
	}
	
	for _, folder := range rootFolderList {
		folderName := folder.Name()
		//check if it's a directory for a thread or not
		if folder.IsDir() && strings.HasPrefix(folderName, "cpu") && !strings.HasSuffix(folderName, "idle") && !strings.HasSuffix(folderName,"freq"){
			coreID, err := os.ReadFile(filepath.Join(rootPath, folderName, topologyFilePath))
			if err != nil {
				return err
			}
			//convert the data (bytes - coreId) to a string -> remove trailing spaces -> convert to an int
			intCoreID, iErr := strconv.Atoi(strings.TrimSpace(string(coreID)))
			if iErr != nil {
				return iErr
			}
			// the folder name has the ID of the thread (cpu0, cpu1 etc.)
			threadID, tErr := strconv.Atoi(strings.Replace(folderName, "cpu", "", 1))
			if tErr != nil {
				return tErr
			}
			//create the topology map. The physycal core is the key
			cr.CPUTopology[intCoreID] = append(cr.CPUTopology[intCoreID], threadID)
		}
	}
	return nil
}

func (cr *CPUReading) mapCPULoadFile() error {
	rootPath := "/proc/stat"
	file, err := os.Open(rootPath)
	if err !=nil {
		return errors.New("Error opening the /proc/stat file")
	}
	cr.StatFile = file
	return nil
}

// GetReady discovers CPU temperature sensors and prepares
// the file descriptors required to read temperature values.
// The files are opened and cached.
func (cr *CPUReading) GetReady() error {
	// creates the map for the raw readings
	cr.RawReadings = make([]CoreInfo, 0, 32)
	// creates the map for CPU Topology
	cr.CPUTopology = make(map[int][]int)
	// creates the map for the reverse CPU Topology
	cr.rCPUTopology = make(map[int]int)


	loadErr := cr.mapCPULoadFile()
	if loadErr != nil {
		return loadErr
	}

	// check cpu topology
	topErr := cr.getCPUTopology()
	if topErr != nil {
		return topErr
	}

	//create a reverse CPU CPUTopology map used by the GetCPULoad(). It's a map[ThreadId]CoreID -> it maps threads to their physical cores
	for physCoreID, threadList := range cr.CPUTopology {
		for _, t := range threadList {
			 cr.rCPUTopology[t] = physCoreID
		}
	}


	//sorted CPU Topology map by CPU Core Number
	keys := make([]int, 0, len(cr.CPUTopology))

	for k := range cr.CPUTopology {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	
	// build the RawReadings as a sorted slice of CORES -> Threads
	for _, core := range keys {
		threads := make([]ThreadInfo, 0, len(cr.CPUTopology[core]))
		for _, thread := range cr.CPUTopology[core] {
			threads = append(threads, ThreadInfo{
				CPU: thread,
				})
		}
		sort.Slice(threads, func(i, j int) bool {
			return threads[i].CPU < threads[j].CPU
		})
		cr.RawReadings = append(cr.RawReadings, CoreInfo{
			CoreID: core,
			Temp: -1,
			Threads: threads,

		})
	}

	//check if any sensor are available. Get the folder path if it is
	rootPath, fErr := cr.findFolder()
	if fErr != nil {
		return fErr
	}
	tErr := cr.mapTempFiles(rootPath)
	if tErr != nil {
		return tErr
	}

	return nil
}

// GetTemp goes over the stored File Descriptors and updates the temperature for each one.
// Mutates the existing struct.
func (cr *CPUReading) GetTemp () error {
	buff := make([]byte, 8)
	errorList := make([]error, 0 , len(cr.RawReadings))
	for i, coreInfo := range cr.RawReadings {
		file := coreInfo.File
		fErr := resetOpenFile(file)
		if fErr != nil {
			errorList = append(errorList, fErr)
			continue
		} 

		n, err := file.Read(buff)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}

		//converts the raw temp(ASCII output) to a string and then an int32
		value, pErr := strconv.Atoi((strings.TrimSpace(string(buff[:n]))))
		if pErr != nil {
			errorList = append(errorList, pErr)
		}
		//stores milidegress as int32
		cr.RawReadings[i].Temp = int32(value)
	}
	var err error
	for _, v := range errorList {
		err = errors.Join(v)
	}
	return err
}

// GetCPULoad get the load per thread and overall load on the cpu.
// It matches this information with each physhical core.
// It also assumes that the sysconf(_SC_CLK_TCK) = 100.
func (cr *CPUReading) GetCPULoad () error {	

	rErr := resetOpenFile(cr.StatFile)
	if rErr != nil {
		return errors.New(fmt.Sprintln("Error while resetting the stat file for the GetCPULoad:\n", rErr))
	}
	scanner := bufio.NewScanner(cr.StatFile)

	//in the /proc/stat file, all CPU's are listed in order. Top being the aggregate one.
	//using this to track which line corresponds to which CPU Thread
	cpuID := -1

	for scanner.Scan() {
		//busyTime = user + nice + system + irq + softirq + steal + guest + guest_nice
		var busyTime uint64 = 0
		//total time = busyTime + iowait + idle
		var totalTime uint64 = 0
		var ok bool

		line := scanner.Bytes()
		busyTime, totalTime, ok = cr.parseStatFile(line)
		if !ok {
			break
		}
		
		//this is the total for the sytem. It's an aggregate for all threads
		if cpuID == -1 {
			//calculates the deltas 
			deltaBusy := float64(busyTime) - float64(cr.TotLoad.prevBusyTime)
			deltaTotal := float64(totalTime) - float64(cr.TotLoad.prevTotTime)

			//caches current readings
			cr.TotLoad.prevTotTime = totalTime
			cr.TotLoad.prevBusyTime = busyTime

			//calculates the actual load
			if deltaTotal > 0 {
				cr.TotLoad.Load = deltaBusy / deltaTotal * 100.0
			}
			cpuID++
			continue
		}

		//using the reverse CPU topology list to map the thread to the physical core
		coreID := cr.rCPUTopology[cpuID]

		//calculates the actual thread load
		//assigns it directly to the correct thread
		for i, v := range cr.RawReadings[coreID].Threads {
			if v.CPU == cpuID {
				//calculates the deltas 
				deltaBusy := float64(busyTime) - float64(cr.RawReadings[coreID].Threads[i].prevBusy)
				deltaTotal := float64(totalTime) - float64(cr.RawReadings[coreID].Threads[i].prevTotal)

				if deltaTotal > 0 {
				cr.RawReadings[coreID].Threads[i].Load = deltaBusy / deltaTotal * 100.0
				
			}
				//caches current readings
				cr.RawReadings[coreID].Threads[i].prevBusy = busyTime
				cr.RawReadings[coreID].Threads[i].prevTotal = totalTime
			}
		}
		cpuID++
	}
	sErr := scanner.Err()
		if sErr != nil {
			return errors.New(fmt.Sprintln("Error parsing stat file on line ", cpuID + 1, ". Error:\n", sErr))
		} 

	return nil
}

//parseStatFile is a helper function to parse each line of the Stat file
//only used by the GetCPULoad
func (cr *CPUReading) parseStatFile (line []byte) (busy, total uint64, ok bool) {
	if line[0] != 'c' || line [1] != 'p' || line [2] != 'u' {
		return 0, 0, false
	}

	lineLen := len(line)
	var (
		i     int = 0
		field int = 0
		val   uint64 
	)
	busy = 0
	total = 0

	// helper: flush current value into correct slot
	flush := func(v uint64, idx int) {
		switch idx {
		case 1: // user
			busy += v
			total += v
		case 2: // nice
			busy += v
			total += v
		case 3: // system
			busy += v
			total += v
		case 4: // idle
			total += v
		case 5: // iowait
			total += v
		case 6: // irq
			busy += v
			total += v
		case 7: // softirq
			busy += v
			total += v
		case 8: // steal
			busy += v
			total += v
		case 9: // guest
			busy += v
			total += v
		case 10: // guest_nice
			busy += v
			total += v
		}
	}
	//skip cpu
	for i < lineLen && line[i] != ' ' {
		i++
	}
	//skip the space
	i++

	//skip the spaces
	for i < lineLen {
		//skip the spaces
		for i < lineLen && line[i] == ' '{
			i++
		}

		if i >= lineLen {
			break
		}

		field++
		val = 0

		//parse number
		for i < lineLen && line[i] >= '0' && line[i] <= '9' {
			//these are ASCII characters (digits) so they are stored as uints
			// '0' is the ASCII character for the uint 48
			// '9' is the ASCII character for the uint 57
			// 57 - 48 = 9
			val = val*10 + uint64(line[i] - '0')
			i++
		}

		flush(val, field)
	}

	ok = field >= 8
	return busy, total, ok

}