package synth

type NoteState int

const (
	Attack NoteState = iota
	Decay
	Sustain
	Release
)

type Note struct {
	pitch int8
	state NoteState
}

// Pitch is counted as an offset (halftones) from the A4
func NewNote(pitch int8, velocity uint8) Note {
	return Note{
		pitch: pitch,
		state: Attack,
	}
}

func (note Note) GetState() NoteState {
	return note.state
}

func (note Note) NextState() Note {
	note.state += 1

	return note
}
