// capture some audio from the default device, and play it back.
//
// to run this example, simply `go run examples/echo.go`.
package main

import (
	".." // pulse-simple
	"encoding/binary"
	"fmt"
	"math"
)

func main() {
	// define the sample spec we want,
	// in this case single-channel 44100Hz little-endian 32-bit float.
	ss := pulse.SampleSpec{pulse.SAMPLE_FLOAT32LE, 44100, 1}

	// create the capture device and audio storage buffers
	capture, err := pulse.Capture("pulse-simple test", "capture test", &ss)
	defer capture.Free()
	if err != nil {
		fmt.Printf("Could not create capture stream: %s\n", err)
		return
	}
	storage := make([]byte, 44100*4)
	// 1/25 second buffer for volume meter and stop button responsiveness
	buffer := make([]byte, 44100*4/25)
	fmt.Println("Recording... (press ENTER to finish capture)")

	// wait in the background for the user to press enter
	stop := make(chan bool)
	go func() {
		var input string
		fmt.Scanln(&input)
		stop <- true
	}()

	// until the user presses enter, continue to capture audio
Capture:
	for {
		select {
		case <-stop:
			break Capture
		default:
			_, err := capture.Read(buffer)
			if err != nil {
				fmt.Printf("Aborting due to audio capture error: %s\n", err)
				break Capture
			}
			storage = append(storage, buffer...)
			// also display a crude volume-meter
			display_volume(buffer)
		}
	}
	fmt.Print("\r") // overwrite the enter from input

	// create playback stream and play back the captured audio
	playback, err := pulse.Playback("pulse-simple test", "playback test", &ss)
	defer playback.Free()
	defer playback.Drain()
	if err != nil {
		fmt.Printf("Could not create playback stream: %s\n", err)
		return
	}
	fmt.Println("Playing back recorded audio...")
	playback.Write(storage)
}

func display_volume(buffer []byte) {
	numsamples := len(buffer) / 4
	sumsquares := float64(0)
	for i := 0; i < numsamples; i++ {
		bits := binary.LittleEndian.Uint32(buffer[4*i : 4*i+4])
		value := math.Float32frombits(bits)
		sumsquares += float64(value * value)
	}
	// root mean square of the signal should give an approximation of volume
	volume := math.Sqrt(sumsquares / float64(numsamples))

	// print between 0 and 40 '#' characters as a volume meter
	maxvolume := 0.5
	fmt.Print("\r")
	for i := 0.0; i < maxvolume; i += maxvolume / 40 {
		if i < volume {
			fmt.Print("#")
		} else {
			fmt.Print(" ")
		}
	}
}
