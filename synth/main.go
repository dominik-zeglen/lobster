package synth

import (
	"fmt"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	a          = 440
	channels   = 16
	chunkSize  = 4096
	sampleRate = 44100
)

func Start() {
	oscs := []Oscillator{
		Oscillator{
			shape: Sine,
		},
		Oscillator{
			shape:  Sine,
			detune: 4,
		},
		Oscillator{
			shape:  Sine,
			detune: 7,
		},
	}

	mixer := NewMixer()
	player := newPlayer()
	defer player.close()

	getAudio := func(ch chan []float64) {
		for oscInd := 0; oscInd < len(oscs); oscInd++ {
			mixer.AddChunk(oscs[oscInd].GetChunk(), oscInd)
		}
		audio := mixer.Mix()

		ch <- audio
	}

	alive := true
	var wg sync.WaitGroup
	wg.Add(3)

	// Play loop
	go func() {
		defer wg.Done()

		audioCh := make(chan []float64)
		go getAudio(audioCh)
		audio := <-audioCh

		for alive {
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
			case 113:
				alive = false

			case 119:
				oscs[0].detune++
			case 115:
				oscs[0].detune--

			case 101:
				oscs[1].detune++
			case 100:
				oscs[1].detune--

			case 114:
				oscs[2].detune++
			case 102:
				oscs[2].detune--
			}
		}
	}()

	// UI loop
	go func() {
		defer wg.Done()

		first := true
		for alive {
			msg := fmt.Sprintf(
				"Osc 1: %d, Osc 2: %d, Osc 3: %d",
				oscs[0].detune,
				oscs[1].detune,
				oscs[2].detune,
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

	wg.Wait()

	// utils.SaveToWav("output.wav", sampleRate, wave)
	// utils.SaveToJson("output.json", wave)

}
