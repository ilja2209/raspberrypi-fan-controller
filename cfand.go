package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const defaltFanPin = 12
const defaultSleepTime = "2s"
const defaultLowTempTreshold = 40.0  //Celsius
const defaultHighTempTreshold = 50.0 //Celsius

func getCPUTemp() float64 {
	cmd := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
	if output, err := cmd.Output(); err != nil {
		panic(err)
	} else {
		var strOutput string = string(output)
		if value, err := strconv.ParseFloat(strings.Trim(strOutput, "\n"), 32); err != nil {
			panic(err)
		} else {
			temp := value / 1000.0
			return temp
		}
	}
}

func switchFan(pinNumber int, isOn bool) {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pin := rpio.Pin(pinNumber)
	pin.Output()
	if isOn {
		pin.High()
	} else {
		pin.Low()
	}
	rpio.Close()
}

func main() {

	if _, err := net.Listen("unix", "@/tmp/cfand"); err != nil {
		fmt.Println("CPU fan controller was already running")
		fmt.Println("CPU temperature is: ", getCPUTemp())
		return
	}

	var lowTempTreshold = *flag.Float64("low-temp-treshold", defaultLowTempTreshold, "This parameter sets low tmperature treshold when the fan is switched off")
	var highTempTreshold = *flag.Float64("high-temp-treshold", defaultHighTempTreshold, "This parameter sets high tmperature treshold when the fan is switched on")
	var pinNumber = *flag.Int("pin", defaltFanPin, "This parameter sets a pin of raspberry pi to which the fan is connected")
	flag.Parse()

	sleepTime, err := time.ParseDuration(defaultSleepTime)
	if err != nil {
		panic(err)
	}
	fmt.Println("CPU fan controller is running")
	for {
		net.Listen("unix", "@/tmp/cfand")
		temp := getCPUTemp()
		if temp < lowTempTreshold {
			switchFan(pinNumber, false)
		}

		if temp > highTempTreshold {
			switchFan(pinNumber, true)
		}
		time.Sleep(sleepTime)
	}
}
