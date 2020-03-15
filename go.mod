module lobster

go 1.12

require (
	github.com/eiannone/keyboard v0.0.0-20190314115158-7169d0afeb4f // indirect
	github.com/go-audio/wav v1.0.0 // indirect
	github.com/hajimehoshi/oto v0.5.4 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
	gitlab.com/gomidi/midi v1.14.1
	gitlab.com/gomidi/rtmididrv v0.4.2
	lobster/synth v0.0.0
	lobster/utils v0.0.0 // indirect
)

replace lobster/synth => ./synth

replace lobster/utils => ./utils
