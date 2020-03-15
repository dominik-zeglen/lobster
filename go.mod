module lobster

go 1.12

require (
	github.com/eiannone/keyboard v0.0.0-20190314115158-7169d0afeb4f // indirect
	github.com/go-audio/audio v1.0.0
	github.com/go-audio/wav v1.0.0
	github.com/hajimehoshi/oto v0.5.4
	github.com/stretchr/testify v1.5.1 // indirect
	lobster/synth v0.0.0
	lobster/utils v0.0.0
)

replace lobster/synth => ./synth

replace lobster/utils => ./utils
