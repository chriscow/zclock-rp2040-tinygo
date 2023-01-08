//go:build rp2040
// +build rp2040

package main

import (
	"fmt"
	"machine"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/colornames"
)

func getUserTime() time.Time {

	println("\r                     S T A H U R R I C A N E")
	println("\r")
	println("                            ####################,\r")
	println("                       #############################\r")
	println("                   ,#######################\r")
	println("                  ############################\r")
	println("                 ##############*###############\r")
	println("                ###########           ###########\r")
	println("               ##########              ##########\r")
	println("               ##########               #########\r")
	println("               /#########/             ##########\r")
	println("                ###########          *###########\r")
	println("                 ###############################\r")
	println("                  ###########################\r")
	println("                    #######################\r")
	println("     ##################################\r")
	println("             *################.\r\n\r")
	println("\r")
	println("=== ZETA CLOCK TIME CONFIG ===")
	println("Enter the current time (HH:MM) and press enter:")

	buffer := make([]byte, 0)

	for {
		if machine.Serial.Buffered() > 0 {
			data, _ := machine.Serial.ReadByte()

			if data == '\r' || data == '\n' {
				break
			}

			buffer = append(buffer, data)
			fmt.Println()
		}
	}

	tokens := strings.Split(string(buffer), ":")
	now := time.Now()
	hour, err := strconv.Atoi(tokens[0])
	if err != nil {
		fmt.Println("\r\n\r\nInvalid hour value.", err)
		return time.Unix(0, 0)
	}
	min, err := strconv.Atoi(tokens[1])
	if err != nil {
		fmt.Println("\r\n\r\nInvalid minute value.", err)
		return time.Unix(0, 0)
	}

	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, time.Local)
}

func main() {
	mcu, err := newDisplay()
	if err != nil {
		for {
			time.Sleep(3 * time.Second)
			fmt.Println(err)
		}
	}

	machine.InitSerial()

	spiral := &spiral{}
	last := time.Unix(0, 0)

	var minHand, hourHand Line
	now := time.Now()
	timeSet := false
	min := minSinceMidnight(now)
	t := timedata[min]
	im := t.imaginary
	mi := t.index

	var offset time.Duration

	for {

		now = time.Now().Add(offset)

		// if the time hasn't been set yet
		// and the user pressed enter, get the time
		if !timeSet && machine.Serial.Buffered() > 0 {
			// if the user pressed enter go into setup
			for machine.Serial.Buffered() > 0 {
				b, _ := machine.Serial.ReadByte()
				if b == '\r' || b == '\n' {
					offset = time.Until(getUserTime())
					now = time.Now().Add(offset)
					last = time.Unix(0, 0)
					timeSet = true
					fmt.Println("time set to:", now, "\r")
					break
				}
			}
		}

		mcu.FillDisplay(colornames.Black)
		spiral.calc(.5, im)

		// Every minute, get the new imaginary value
		if now.Sub(last) > time.Minute {

			min = minSinceMidnight(now)
			t = timedata[min]
			im = t.imaginary
			mi = t.index

			spiral.calc(.5, im)
			minHand.A = spiral.joints[mi]
			minHand.B = spiral.joints[mi+1]
			hourHand.A = spiral.joints[mi+1]
			hourHand.B = spiral.joints[mi+2]

			last = now
		}
		drawSpiral(mcu, spiral, mi)
		if timeSet {
			drawHands(mcu, minHand, hourHand)
		}
		mcu.Display()

		im += .004
	}
}
