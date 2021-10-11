package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ilja2209/cfand/pidctrl"

	"github.com/stianeikeland/go-rpio"
)

const defaltFanPin = 12
const defaultSleepTime = "1s"
const defaultLowTempTreshold = 40.0     // Celsius
const defaultHighTempTreshold = 50.0    // Celsius
const defaultSetPointTemperature = 45.0 // Celsius
const maxSpeed = 255.0
const baseFreq = 38000 // Hz
const minOutputLimit = -10
const maxOutputLimit = 10
const pCoef = 0.5
const iCoef = 0.5
const dCoef = 0.5

// For raspberry PI 4 only
var pinsWithPWM = map[int]bool{
	12: true,
	13: true,
	18: true,
	19: true,
	40: true,
	41: true,
	45: true,
}

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

func getPin(pinNumber int) rpio.Pin {
	if err := rpio.Open(); err != nil {
		panic(err)
	}

	return rpio.Pin(pinNumber)
}

func switchFan(pinNumber int, isOn bool) {
	pin := getPin(pinNumber)
	pin.Output()
	if isOn {
		pin.High()
	} else {
		pin.Low()
	}
	rpio.Close()
}

func setFanSpeed(pinNumber int, speed float64) {
	// fmt.Println("  V = ", speed)
	pin := getPin(pinNumber)
	pin.Pwm()
	pin.DutyCycle(uint32(speed), uint32(maxSpeed))
	pin.Freq(baseFreq * int(maxSpeed))
}

func fanPulseController(currentTemp float64, lowTempTreshold float64, highTempTreshold float64, pinNumber int) {
	if currentTemp < lowTempTreshold {
		switchFan(pinNumber, false)
	}

	if currentTemp > highTempTreshold {
		switchFan(pinNumber, true)
	}
}

func initPIDController(setTemp float64) *pidctrl.PIDController {
	pidController := pidctrl.NewPIDController(pCoef, iCoef, dCoef)
	pidController.Set(setTemp)
	pidController.SetOutputLimits(minOutputLimit, maxOutputLimit)
	return pidController
}

func fanPIDController(pidController *pidctrl.PIDController, currentTemp float64, pinNumber int) {
	y := pidController.UpdateDuration(currentTemp, time.Second)
	speed := -(y - 10) * 12.75 // scale [-10; 10] to [255; 0]
	setFanSpeed(pinNumber, speed)
}

func main() {
	fmt.Println("CPU temperature is: ", getCPUTemp())

	var setPointTemperature = flag.Float64("setpoint-temperature", defaultSetPointTemperature, "This parameter sets target value of CPU temperature")
	var lowTempTreshold = flag.Float64("low-temp-threshold", defaultLowTempTreshold, "This parameter sets low temperature threshold when the fan is switched off")
	var highTempTreshold = flag.Float64("high-temp-threshold", defaultHighTempTreshold, "This parameter sets high temperature threshold when the fan is switched on")
	var pinNumber = flag.Int("pin", defaltFanPin, "This parameter sets a pin of raspberry pi to which the fan is connected")
	flag.Parse()

	sleepTime, err := time.ParseDuration(defaultSleepTime)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Setpoint", *setPointTemperature)

	pidController := initPIDController(*setPointTemperature)

	isPID := pinsWithPWM[*pinNumber]

	fmt.Println("CPU fan controller is running")

	for {
		temp := getCPUTemp()
		// fmt.Print("T = ", temp)

		if isPID {
			fanPIDController(pidController, temp, *pinNumber)
		} else {
			fanPulseController(temp, *lowTempTreshold, *highTempTreshold, *pinNumber)
		}

		time.Sleep(sleepTime)
	}
}
