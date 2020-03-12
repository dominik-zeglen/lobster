module lobster

go 1.12

require (
	github.com/go-audio/audio v1.0.0
	github.com/go-audio/wav v1.0.0
	lobster/synth v0.0.0
	lobster/utils v0.0.0
)

replace lobster/synth => ./synth

replace lobster/utils => ./utils
