package synth

import "io"

type AudioBus struct {
	alive  *bool
	notes  []Note
	oscs   []Oscillator
	volume *int
}

func (bus *AudioBus) RegisterNote(note Note) {
	bus.notes = append(bus.notes, note)
}

func (bus AudioBus) Read(buf []byte) (int, error) {
	if !*bus.alive {
		return 0, io.EOF
	}
	out := make([]int32, chunkSize)

	for noteInd := range bus.notes {
		chunks := make([][]int16, len(bus.oscs))
		for oscInd := range bus.oscs {
			chunks[oscInd] = bus.oscs[oscInd].GetChunk(bus.notes[noteInd].pitch)
		}

		for i := 0; i < chunkSize; i++ {
			for oscInd := range bus.oscs {
				out[i] += int32(chunks[oscInd][i])
			}
		}
	}

	for sampleInd := range out {
		s := int16(out[sampleInd] / int32(len(bus.notes)) / int32(len(bus.oscs)) * int32(*bus.volume) / 100)
		buf[sampleInd*2] = byte(s)
		buf[sampleInd*2+1] = byte(s >> 8)
	}

	return chunkSize * 2, nil
}
