# Zeta Clock

A unique timepiece for RP2040-based round LCD displays that visualizes time through an animated Riemann Zeta spiral. Built specifically for the Waveshare RP2040 LCD 1.28" development board.

## Features
- Dynamic spiral animation that morphs based on the current time
- Hour and minute hands represented by colored spiral segments
- Motion-based sleep mode using QMI8658C IMU
- Battery voltage monitoring with visual indicator
- Interactive time setting via serial interface

## Hardware Requirements
- [Waveshare RP2040-LCD-1.28](https://www.waveshare.com/wiki/RP2040-LCD-1.28) development board
- USB cable for programming

## Building and Flashing

1. Install TinyGo following the [official instructions](https://tinygo.org/getting-started/)

2. Clone this repository:
```bash
git clone https://github.com/yourusername/zeta-clock.git
cd zeta-clock
```

3. Flash to your device:
```bash
tinygo flash -target=pico .
```

## Usage
1. Connect to the device's serial port (115200 baud)
2. Enter current time in HH:MM format when prompted
3. The spiral will animate and show the current time:
   - Orange segment: minute hand
   - Red segment: hour hand
   - Gray spiral: time reference points

## Credits
Vector math utilities adapted from [Pixel Game Engine](https://github.com/faiface/pixel)

## License
MIT License
