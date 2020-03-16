package synth

import (
	"fmt"
	"io"
	"lobster/ioMidi"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/hajimehoshi/oto"
)

const (
	a          = 440
	channels   = 16
	chunkSize  = 2048
	sampleRate = 44100
)

func Start() {
	volume := 10

	oscs := []Oscillator{
		Oscillator{
			shape: Sine,
		},
		Oscillator{
			detune: 4,
			shape:  Sine,
		},
		Oscillator{
			detune: 7,
			shape:  Sine,
		},
		Oscillator{
			detune: 10,
			shape:  Sine,
		},
	}

	alive := true
	var wg sync.WaitGroup
	wg.Add(4)

	audioBus := &AudioBus{
		alive:  &alive,
		oscs:   oscs,
		volume: &volume,
	}

	addNote := func(pitch int8) {
		audioBus.RegisterNote(Note{
			pitch: pitch,
		})
	}

	removeNote := func(pitch int8) {
		audioBus.UnregisterNote(Note{
			pitch: pitch,
		})
	}

	// Play loop
	go func() {
		defer wg.Done()

		audioCtx, err := oto.NewContext(sampleRate, 1, 2, 2)
		if err != nil {
			panic(err)
		}
		defer audioCtx.Close()
		player := audioCtx.NewPlayer()

		// for alive {

		// }
		if _, err := io.Copy(player, audioBus); err != nil {
			panic(err)
		}
		if err := player.Close(); err != nil {
			panic(err)
		}
	}()

	// Input loop
	go func() {
		defer wg.Done()
		for alive {
			char, _, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}
			switch char {
			case 'q':
				alive = false

			case 'w':
				oscs[0].detune++
			case 's':
				oscs[0].detune--

			case 'e':
				oscs[1].detune++
			case 'd':
				oscs[1].detune--

			case 'r':
				oscs[2].detune++
			case 'f':
				oscs[2].detune--

			case 't':
				oscs[3].detune++
			case 'g':
				oscs[3].detune--

			case '=':
				if volume < 100 {
					volume++
				}
			case '-':
				if volume > 0 {
					volume--
				}
			}
		}
	}()

	// UI loop
	go func() {
		defer wg.Done()

		first := true
		for alive {
			msg := fmt.Sprintf(
				"Volume %d, Osc 1: %d, Osc 2: %d, Osc 3: %d, Osc 4: %d",
				volume,
				oscs[0].detune,
				oscs[1].detune,
				oscs[2].detune,
				oscs[3].detune,
			)
			if !first {
				cls := ""
				for i := 0; i < len(msg); i++ {
					cls = cls + "\b"
				}
				fmt.Print(cls)
			}
			fmt.Print(msg)
			time.Sleep(time.Millisecond * 100)
			first = false
		}
	}()

	// Midi IO loop
	go ioMidi.Loop(&alive, &wg, addNote, removeNote)

	wg.Wait()

	// utils.SaveToWav("output.wav", sampleRate, wave)
	// utils.SaveToJson("output.json", wave)

}
