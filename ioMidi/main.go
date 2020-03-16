package ioMidi

import (
	"fmt"
	"sync"

	"gitlab.com/gomidi/midi/mid"
	driver "gitlab.com/gomidi/rtmididrv"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func printPort(port mid.Port) {
	fmt.Printf("[%v] %s\n", port.Number(), port.String())
}

func printInPorts(ports []mid.In) {
	fmt.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

func printOutPorts(ports []mid.Out) {
	fmt.Printf("MIDI OUT Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

type MidiIoOpts struct {
	Pitch *int8
}

func Loop(alive *bool, wg *sync.WaitGroup, addNote func(int8), removeNote func(int8)) {
	defer wg.Done()

	noteOn := func(p *mid.Position, channel, key, vel uint8) {
		addNote(int8(key - 60))
	}

	noteOff := func(p *mid.Position, channel, key, vel uint8) {
		removeNote(int8(key - 60))
	}

	drv, err := driver.New()
	must(err)

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	// printInPorts(ins)
	// printOutPorts(outs)

	must(ins[1].Open())

	rd := mid.NewReader()

	rd.Msg.Channel.NoteOn = noteOn
	rd.Msg.Channel.NoteOff = noteOff
	rd.Msg.Each = nil
	// rd.Msg.Channel.NoteOff = noteOff

	var audioWg sync.WaitGroup
	audioWg.Add(1)
	// listen for MIDI
	go mid.ConnectIn(ins[1], rd)
	// mid.
	audioWg.Wait()

	// for *alive {
	// 	b, err := in.SetListener()
	// }
}
