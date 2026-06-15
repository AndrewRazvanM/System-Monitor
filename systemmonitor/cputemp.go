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
			fmt.Println("Found Folder:", filePath)
			return filePath, nil
		}
		}
	}
	return "", errors.New("Could not find a sensor")
}

func resetOpenFile (file *os.File) error {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
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

// GetReady discovers CPU temperature sensors and prepares
// the file descriptors required to read temperature values.
// The files are opened and cached.
func (cr *CPUReading) GetReady() error {
	// create the maps
	cr.FilesDescriptor = make(map[string]*os.File)
	cr.RawTemps = make(map[string]int32)

	//check if any sensor are available. Get the folder path if it is
	rootPath, err := cr.findFolder()
	if err != nil {
		return err
	}

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
			label := strings.TrimSpace(string(data))

			//check for the input file -> it contains the raw temp 
			inputStr := strings.TrimSuffix(fileName, "label") + "input"
			inputPath := filepath.Join(rootPath, inputStr)
			FileDescriptor, iError := os.Open(inputPath)
			if iError != nil {
				return iError
			}
			cr.FilesDescriptor[label] = FileDescriptor
		}
	}
	return nil
}

// GetTemp goes over the stored File Descriptors and updates the temperature for each one.
// Mutates the existing struct.
func (cr *CPUReading) GetTemp () error {
	buff := make([]byte, 8)
	errorList := make([]error, 0 , len(cr.FilesDescriptor) - 1)
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
		cr.RawTemps[core] = int32(value)
	}
	var err error
	for _, v := range errorList {
		err = errors.Join(v)
	}
	return err
}

