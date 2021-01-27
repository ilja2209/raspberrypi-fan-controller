# Fan controller for Raspberry Pi 4 based on PID regulator
The fan controller which can be run on your raspberry pi to control CPU temperature. 
It can use PID regulator to ajust CPU temperature more precisely or just pulse regulator (depends on pins the fan is connected which). 
PID regulator can dramatically reduce a noise of the CPU fan (Depends on settings).
To use PID regulator you need to connect your fan to any of 12, 13, 18, 19, 40, 41, 45 pins and set target temperature using parameter "setpoint-temperature" 

The application can be run with several parameters:
```
-pin                 This parameter sets a pin of raspberry pi to which the fan is connected

For PID regulator:
-setpoint-temperature This parameter sets target value of CPU temperature. (Deault value is 45 Celsius)

For pulse regulator: 
-low-temp-treshold   This parameter sets low tmperature treshold when the fan is switched off. (Deault value is 40 Celsius)
-high-temp-treshold  This parameter sets high tmperature treshold when the fan is switched on. (Deault value is 50 Celsius)
```

To download the application and add to autorun follow the next steps (for Ubuntu):
```
1. cd /etc/init.d
2. sudo wget https://raw.githubusercontent.com/ilja2209/raspberrypi-fan-controller/main/etc/init.d/cfand.sh
3. sudo wget https://github.com/ilja2209/raspberrypi-fan-controller/releases/download/1.0-linux-arm64/cfand
4. Modify cfand.sh to set needed parameters (sudo vim cfand.sh)
5. Make link: sudo ln -s /etc/init.d/cfand /etc/rc2.d/S99cfand
6. Reboot your raspberry pi
7. Enjoy
```

The fan can be connected to raspberry pi 4 using the next circuit:

![alt text](https://github.com/ilja2209/raspberrypi-fan-controller/raw/main/Fan%20controller%20circuit.png?raw=true)
