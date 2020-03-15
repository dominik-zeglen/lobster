package utils

import (
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func floatPcmToInt(pcm []float64) []int {
	intPcm := make([]int, len(pcm))

	for probInd := 0; probInd < len(pcm); probInd++ {
		intPcm[probInd] = int(pcm[probInd] * float64(32768))
	}

	return intPcm
}

func newAudioBuffer(pcm []float64) *audio.IntBuffer {
	buf := audio.IntBuffer{
		Data: floatPcmToInt(pcm),
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  8000,
		},
	}

	return &buf
}

func SaveToWav(fpath string, sampleRate int, pcm []float64) error {
	out, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer out.Close()

	e := wav.NewEncoder(out, sampleRate, 16, 1, 1)

	audioBuf := newAudioBuffer(pcm)

	// Write buffer to output file. This writes a RIFF header and the PCM chunks from the audio.IntBuffer.
	if err := e.Write(audioBuf); err != nil {
		return err
	}
	if err := e.Close(); err != nil {
		return err
	}

	return nil
}
