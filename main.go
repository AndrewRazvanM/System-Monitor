package main

import (
	"fmt"

	"github.com/AndrewRazvanM/System-Monitor/systemmonitor"
)

func main() {
	tempRead:= systemmonitor.CPUReading{}
	err := tempRead.GetReady()
	fmt.Println(tempRead.FilesDescriptor, err)

	rErr := tempRead.GetTemp()
	if rErr != nil {
		fmt.Println("Error getting temp readings: ", rErr)
	}

	errList := tempRead.Close()
	fmt.Println("Errors: ", errList)

	tempRead.FormatReadings()
	fmt.Println("Printing CPU cores temperatures:")
	for core, v := range tempRead.FormattedTemp {
		fmt.Println(core, "\nTemp:", v)
	}
}