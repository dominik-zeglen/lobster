package synth

import (
	"fmt"
	"lobster/ioMidi"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	a          = 440
	channels   = 16
	chunkSize  = 2048
	sampleRate = 44100
)

func Start() {
	volume := 10
	pitch := int8(0)

	oscs := []Oscillator{
		Oscillator{
			pitch: &pitch,
			shape: Sine,
		},
		Oscillator{
			detune: 4,
			pitch:  &pitch,
			shape:  Sine,
		},
		Oscillator{
			detune: 7,
			pitch:  &pitch,
			shape:  Sine,
		},
		Oscillator{
			detune: 10,
			pitch:  &pitch,
			shape:  Sine,
		},
	}

	mixer := NewMixer()
	player := newPlayer()
	defer player.close()

	getAudio := func(ch chan []float64) {
		for oscInd := 0; oscInd < len(oscs); oscInd++ {
			mixer.AddChunk(oscs[oscInd].GetChunk(), oscInd)
		}
		audio := mixer.Mix(volume)

		ch <- audio
	}

	setPitch := func(newPitch int8) {
		pitch = newPitch
	}

	alive := true
	var wg sync.WaitGroup
	wg.Add(4)

	// Play loop
	go func() {
		defer wg.Done()

		audioCh := make(chan []float64)
		go getAudio(audioCh)
		audio := <-audioCh

		for alive {
			// TODO: Instead of using dirty hacks just make oscillator return
			// io.Reader so it won't be rewritten again and again, just create
			// one stream which will act as one note channel
			go getAudio(audioCh)
			err := player.playPcm(audio)
			if err != nil {
				panic(err)
			}

			audio = <-audioCh
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
	go ioMidi.Loop(&alive, &wg, setPitch)

	wg.Wait()

	// utils.SaveToWav("output.wav", sampleRate, wave)
	// utils.SaveToJson("output.json", wave)

}
