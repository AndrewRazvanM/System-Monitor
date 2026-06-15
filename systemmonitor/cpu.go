package systemmonitor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
func (cr *CPUReading) Close() []error {
	err := []error{}
	for _, v := range cr.FilesDescriptor {
		error := v.Close()
		if error != nil {
			errorMsg := fmt.Sprintln("File ", v.Name(), " could not be closed. Error: ", error)
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
			//convert the data from bytes of ASCII number to int
			strData := strings.TrimSpace(string(data))
			//find the inde of the last space in the string -> create a slice from the string,
			//starting from index + 1 -> convert this to an int64
			label, intErr := strconv.ParseInt(strData[strings.LastIndex(strData, " ") + 1:], 10, 32)
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
			cr.FilesDescriptor[int(label)] = FileDescriptor
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

// GetReady discovers CPU temperature sensors and prepares
// the file descriptors required to read temperature values.
// The files are opened and cached.
func (cr *CPUReading) GetReady() error {
	// creates the map for file descriptors
	cr.FilesDescriptor = make(map[int]*os.File)
	// creates the map for the raw readings
	cr.RawReadings = make([]CoreInfo, 0, 32)
	// creates the map for CPU Topography
	cr.CPUTopology = make(map[int][]int)

	//check if any sensor are available. Get the folder path if it is
	rootPath, fErr := cr.findFolder()
	if fErr != nil {
		return fErr
	}
	tErr := cr.mapTempFiles(rootPath)
	if tErr != nil {
		return tErr
	}
	topErr := cr.getCPUTopology()
	if topErr != nil {
		return topErr
	}
	return nil
}

// GetTemp goes over the stored File Descriptors and updates the temperature for each one.
// Mutates the existing struct.
func (cr *CPUReading) GetTemp () error {
	buff := make([]byte, 8)
	errorList := make([]error, 0 , len(cr.FilesDescriptor))
	for core, file := range cr.FilesDescriptor {
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
		value, pErr := strconv.ParseInt(strings.TrimSpace(string(buff[:n])), 10, 32)
		if pErr != nil {
			errorList = append(errorList, pErr)
		}
		//stores milidegress as int32
		cr.RawReadings = append(cr.RawReadings, CoreInfo{CoreID: core, Temp: int32(value)})
	}
	var err error
	for _, v := range errorList {
		err = errors.Join(v)
	}
	return err
}

