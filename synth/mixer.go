package synth

type Mixer struct {
	chunks [][]float64
}

func (mixer *Mixer) reset() {
	mixer.chunks = make([][]float64, channels)
	for chunkInd := range mixer.chunks {
		mixer.chunks[chunkInd] = make([]float64, chunkSize)
	}
}

func NewMixer() Mixer {
	m := Mixer{}
	m.reset()

	return m
}

func (mixer *Mixer) AddChunk(chunk []float64, channel int) {
	mixer.chunks[channel] = chunk
}

func (mixer *Mixer) Mix() []float64 {
	defer mixer.reset()

	mixed := make([]float64, chunkSize)

	for probInd := range mixed {
		var v float64

		for channelInd := range mixer.chunks {
			v += mixer.chunks[channelInd][probInd]
		}

		mixed[probInd] = v / channels
	}

	return mixed
}
