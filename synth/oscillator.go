package synth

import (
	"math"
)

type WaveShape int

const (
	Sine WaveShape = iota
	Triangle
	Square
	Saw
	Noise
)

type Oscillator struct {
	detune      int8
	phaseOffset float64
	shape       WaveShape
}

func (osc *Oscillator) GetChunk() []float64 {
	chunkByRate := float64(chunkSize) / float64(sampleRate)
	wave := make([]float64, chunkSize)
	freq := (a * math.Pow(2, float64(osc.detune)/12))

	for probInd := range wave {
		t := float64(probInd) / sampleRate
		wave[probInd] = math.Sin(2 * math.Pi * (freq*t + osc.phaseOffset))
	}

	_, osc.phaseOffset = math.Modf(freq*chunkByRate + osc.phaseOffset)

	return wave
}
