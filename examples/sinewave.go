// generate a sine wave and send it to the default audio output
package main

import (
	".." // pulse-simple
	"encoding/binary"
	"fmt"
	"math"
)

func main() {
	ss := pulse.SampleSpec{pulse.SAMPLE_FLOAT32LE, 44100, 1}
	pb, err := pulse.Playback("pulse-simple test", "playback test", &ss)
	defer pb.Free()
	defer pb.Drain()
	if err != nil {
		fmt.Printf("Could not create playback stream: %s\n", err)
		return
	}
	playsine(pb, &ss)
}

func playsine(s *pulse.Stream, ss *pulse.SampleSpec) {
	num_notes := 5
	f := []float64{220, 247, 277, 294, 330}
	n := []string{"A3", "B3", "C#4", "D4", "E4"}
	r := float64(ss.Rate)
	data := make([]byte, 4*ss.Rate)
	tau := 2 * math.Pi
	for j := 0; j < num_notes; j++ {
		fmt.Printf("%v\n", n[j])
		for i := 0; i < int(ss.Rate); i++ {
			// (f) Hz sine wave, with 0.5Hz sine envelope over 1 second duration
			sample := float32((math.Sin(tau*f[j]*float64(i)/r) / 3.0) *
				math.Sin((tau/2.0)*float64(i)/r))
			bits := math.Float32bits(sample)
			binary.LittleEndian.PutUint32(data[4*i:4*i+4], bits)
		}
		s.Write(data)
	}
}
