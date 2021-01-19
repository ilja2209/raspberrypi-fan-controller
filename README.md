# Fan controller for Raspberry Pi 4 with PID controller support
The sipliest fan controller which can be run on your raspberry pi to control CPU temperature. It uses simple logic:

The application can be run with several parameters:
```
-low-temp-treshold   This parameter sets low tmperature treshold when the fan is switched off
-high-temp-treshold  This parameter sets high tmperature treshold when the fan is switched on
-pin                 This parameter sets a pin of raspberry pi to which the fan is connected
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
