package main

import (
	".." // pulse-simple
	"fmt"
)

func main() {
	fmt.Printf("Channel Map Test\n")
	fmt.Printf("================\n")
	fmt.Printf("CHANNELS_MAX: %v\n", pulse.CHANNELS_MAX)

	fmt.Printf("\nMono\n")
	fmt.Printf("----\n")
	mono := &pulse.ChannelMap{}
	mono.InitMono()
	print_info(mono)

	fmt.Printf("\nStereo\n")
	fmt.Printf("------\n")
	stereo := &pulse.ChannelMap{}
	stereo.InitStereo()
	print_info(stereo)

	spec := &pulse.SampleSpec{pulse.SAMPLE_S16LE, 44100, 2}
	fmt.Printf("\nspec := &SampleSpec{SAMPLE_S16LE, 44100, 2}\n")
	fmt.Printf("mono.Compatible(spec): %v\n", mono.Compatible(spec))
	fmt.Printf("stereo.Compatible(spec): %v\n", stereo.Compatible(spec))

	fmt.Printf("\n9 Channel AIFF (should fail)\n")
	fmt.Printf("----------------------------\n")
	cmap := &pulse.ChannelMap{}
	err := cmap.InitAuto(9, pulse.CHANNEL_MAP_AIFF)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		print_info(cmap)
	}

	fmt.Printf("\n9 Channel AIFF (should succeed)\n")
	fmt.Printf("-------------------------------\n")
	cmap.InitExtend(9, pulse.CHANNEL_MAP_AIFF)
	print_info(cmap)

	fmt.Printf("\naiff9.Superset(stereo): %v\n", cmap.Superset(stereo))
	fmt.Printf("aiff9.Superset(mono): %v\n", cmap.Superset(mono))

	name, err := stereo.Name()
	if err != nil {
		fmt.Printf("\nstereo.Name(): %v\n", err)
	} else {
		fmt.Printf("\nstereo.Name(): %v\n", name)
	}
	name, err = stereo.PrettyName()
	if err != nil {
		fmt.Printf("stereo.PrettyName(): %v\n", err)
	} else {
		fmt.Printf("stereo.PrettyName(): %v\n", name)
	}

	fmt.Printf("\n7.1 Surround\n")
	fmt.Printf("------------\n")
	err = cmap.InitAuto(8, pulse.CHANNEL_MAP_ALSA)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		print_info(cmap)
		name, err = cmap.PrettyName()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("PrettyName: %v\n", name)
		}
	}

	fmt.Printf("\nsurround has subwoofer: %v\n",
		cmap.HasPosition(pulse.CHANNEL_POSITION_SUBWOOFER))
	fmt.Printf("stereo has subwoofer: %v\n",
		stereo.HasPosition(pulse.CHANNEL_POSITION_SUBWOOFER))

	fmt.Printf("\nsurround.String(): %v\n", cmap.String())
}

func print_info(cmap *pulse.ChannelMap) {
	fmt.Printf("Channels: %v\n", cmap.Channels)
	for i := 0; i < int(cmap.Channels); i++ {
		fmt.Printf("Map[%d]: %d (%v)\n", i, cmap.Map[i], cmap.Map[i])
	}
}
