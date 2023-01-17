//go:build rp2040
// +build rp2040

package main

import (
	"fmt"
	"machine"
	"time"

	"golang.org/x/image/colornames"
)

func main() {
	// mcu, _ := newMCU()

	// mcu.lcd.FillScreen(colornames.Black)

	// tinydraw.Rectangle(mcu, 40, 40, 160, 160, colornames.White)

	// from := V(127, 127)
	// to := V(60, 0).Rotated(math32.Pi / 6).Add(from)
	// tinydraw.Line(mcu, int16(from.X), int16(from.Y), int16(to.X), int16(to.Y), colornames.White)

	// temp := to.Sub(from).Normal().Add(from)
	// tinydraw.Line(mcu, int16(from.X), int16(from.Y), int16(temp.X), int16(temp.Y), colornames.White)

	// mcu.Display()

	// for {
	// 	time.Sleep(time.Second)
	// }

	if err := run(); err != nil {
		for {
			fmt.Println("error calling newDisplay:", err)
			time.Sleep(time.Second)
		}
	}
}

func run() error {
	mcu, err := newMCU()
	if err != nil {
		return err
	}

	spiral := &spiral{}
	everyMin := time.Unix(0, 0)
	everySec := time.Unix(0, 0)

	var minHand, hourHand Line
	now := time.Now()
	timeSet := false
	min := minSinceMidnight(now)
	t := timedata[min]
	im := t.imaginary
	mi := t.index

	var offset time.Duration

	const LCD_SLEEP_ON = 0x10
	const LCD_SLEEP_OFF = 0x11

	for {

		now = time.Now().Add(offset)

		// sleeper.update()
		// sleeping := time.Since(sleeper.lastMovement) > 5*time.Minute
		// deepsleep := time.Since(sleeper.lastMovement) > 15*time.Minute
		// if sleeping {
		// 	mcu.lcd.Command(LCD_SLEEP_ON) // LCD sleep
		// 	mcu.lcd.EnableBacklight(false)

		// 	if deepsleep {
		// 		time.Sleep(2 * time.Second)
		// 	} else {
		// 		time.Sleep(time.Millisecond * 10)
		// 	}

		// 	min = minSinceMidnight(now)
		// 	t = timedata[min]
		// 	im = t.imaginary
		// 	mi = t.index

		// 	continue
		// }

		mcu.lcd.Command(LCD_SLEEP_OFF)
		mcu.lcd.EnableBacklight(true)

		// if the time hasn't been set yet
		// and the user pressed enter, get the time
		if !timeSet && machine.Serial.Buffered() > 0 {
			// if the user pressed enter go into setup
			for machine.Serial.Buffered() > 0 {
				b, _ := machine.Serial.ReadByte()
				if b == '\r' || b == '\n' {
					offset = time.Until(getUserTime())
					now = time.Now().Add(offset)
					everyMin = time.Unix(0, 0)
					timeSet = true
					fmt.Println("time set to:", now, "\r")
					break
				}
			}
		}

		mcu.FillDisplay(colornames.Black)
		spiral.calc(.5, im)

		if now.Sub(everySec) > time.Second {
			// fmt.Printf("bat(v): %.2f\r\n", mcu.Volts())

			everySec = now
		}

		// Every minute, get the new imaginary value
		if now.Sub(everyMin) > time.Minute {

			min = minSinceMidnight(now)
			t = timedata[min]
			im = t.imaginary
			mi = t.index

			spiral.calc(.5, im)
			minHand.B = spiral.joints[mi]
			minHand.A = spiral.joints[mi+1]
			hourHand.A = spiral.joints[mi+1]
			hourHand.B = spiral.joints[mi+2]

			everyMin = now
		}

		spiral.draw(mcu, mi)

		if timeSet {
			drawHand(mcu, minHand, .9, colornames.Orange)
			drawHand(mcu, hourHand, .6, colornames.Red)
		}

		mcu.DrawBattery()
		mcu.Display()

		im += .004
	}
}
