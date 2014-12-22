package main

import (
	".." // pulse-simple
	"fmt"
)

func main() {
	fmt.Printf("spec1 := SampleSpec{SAMPLE_S16LE,44100,2}\n")
	spec1 := create_spec(pulse.SAMPLE_S16LE, 44100, 2)
	print_info(spec1)

	fmt.Printf("\nspec2 := SampleSpec{SAMPLE_FLOAT32,96000,6}\n")
	spec2 := create_spec(pulse.SAMPLE_FLOAT32, 96000, 6)
	print_info(spec2)

	fmt.Printf("\nspec1.Equal(spec2): %v\n", spec1.Equal(spec2))
	fmt.Printf("spec3 := SampleSpec{SAMPLE_S16LE,44100,2}\n")
	spec3 := create_spec(pulse.SAMPLE_S16LE, 44100, 2)
	fmt.Printf("spec3.Equal(spec1): %v\n", spec3.Equal(spec1))

	fmt.Printf("\nspec4 := SampleSpec{SAMPLE_U8,5001,1}\n")
	spec4 := create_spec(pulse.SAMPLE_U8, 5001, 1)
	print_info(spec4)
}

func create_spec(f pulse.SampleFormat, r uint32, c uint8) *pulse.SampleSpec {
	return &pulse.SampleSpec{f, r, c}
}

func print_info(spec *pulse.SampleSpec) {
	fmt.Printf("String:            %v\n", spec.String())
	fmt.Printf("Valid:             %v\n", spec.Valid())
	fmt.Printf("BytesPerSecond:    %v\n", spec.BytesPerSecond())
	fmt.Printf("FrameSize:         %v\n", spec.FrameSize())
	fmt.Printf("SampleSize:        %v\n", spec.SampleSize())
	fmt.Printf("BytesToUsec(64):   %v\n", spec.BytesToUsec(64))
	fmt.Printf("UsecToBytes(100):  %v\n", spec.UsecToBytes(100))
	fmt.Printf("SampleFormat:      %s\n", spec.Format.String())
	fmt.Printf("Format SampleSize: %v\n", spec.Format.SampleSize())
	fmt.Printf("Little-endian:     %v\n", spec.Format.IsLe())
	fmt.Printf("Big-endian:        %v\n", spec.Format.IsBe())
	fmt.Printf("Native-endian:     %v\n", spec.Format.IsNe())
	fmt.Printf("Reverse-endian:    %v\n", spec.Format.IsRe())
}
