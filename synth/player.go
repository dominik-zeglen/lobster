package synth

import (
	"io"

	"github.com/hajimehoshi/oto"
)

type Player struct {
	audioCtx oto.Context
	p        *oto.Player
}

func newPlayer() Player {
	audioCtx, err := oto.NewContext(
		sampleRate,
		1,
		2,
		2,
	)

	if err != nil {
		panic(err)
	}

	player := audioCtx.NewPlayer()

	return Player{
		audioCtx: *audioCtx,
		p:        player,
	}
}

type PcmBuffer struct {
	data      []float64
	length    int64
	pos       int64
	remaining []byte
}

func pcmAsBuffer(pcm []float64) *PcmBuffer {
	const max = 32767

	return &PcmBuffer{
		data:   pcm,
		length: int64(len(pcm)),
	}
}

func (pcmBuf *PcmBuffer) Read(buf []byte) (int, error) {
	if len(pcmBuf.remaining) > 0 {
		n := copy(buf, pcmBuf.remaining)
		pcmBuf.remaining = pcmBuf.remaining[n:]
		return n, nil
	}

	if pcmBuf.pos == pcmBuf.length {
		return 0, io.EOF
	}

	eof := false
	if pcmBuf.pos+int64(len(buf)) > pcmBuf.length {
		buf = buf[:pcmBuf.length-pcmBuf.pos]
		eof = true
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	p := pcmBuf.pos / 2

	for i := 0; i < len(buf)/2; i++ {
		const max = 32767
		b := int16(pcmBuf.data[i] * max)
		buf[2*i] = byte(b)
		buf[2*i+1] = byte(b >> 8)
		p++
	}

	pcmBuf.pos += int64(len(buf))

	n := len(buf)
	if origBuf != nil {
		n = copy(origBuf, buf)
		pcmBuf.remaining = buf[n:]
	}

	if eof {
		return n, io.EOF
	}
	return n, nil
}

func (player Player) close() error {
	return player.audioCtx.Close()
}

func (player Player) playPcm(pcm []float64) error {
	buf := pcmAsBuffer(pcm)

	if _, err := io.Copy(player.p, buf); err != nil {
		return err
	}
	if err := player.p.Close(); err != nil {
		return err
	}
	return nil
}
