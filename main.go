package main

import (
	"math"
	"time"

	"github.com/fr3fou/gusic/gusic"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Key struct {
	gusic.Note
	IsSemitone bool
}

func NewKey(note gusic.Note, isSemitone bool) Key {
	return Key{Note: note, IsSemitone: isSemitone}
}

func (p *Key) Samples() []float32 {
	return samplesToFloat32(
		p.Note.Samples(
			// TODO, configurable params
			48000,
			math.Sin,
			gusic.NewLinearADSR(
				gusic.NewRatios(0.25, 0.25, 0.25, 0.25), 1.35, 0.35,
			),
		),
	)
}

func main() {
	width := int32(1680)
	height := int32(900)
	rl.InitWindow(width, height, "goda - a simple music pad")

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	stream := rl.InitAudioStream(48000, 32, 1)
	defer rl.CloseAudioStream(stream)

	// maxSamples := 48000 * 5
	// maxSamplesPerUpdate := 4096

	// data := make([]float32, maxSamples)

	rl.PlayAudioStream(stream)

	// totalSamples := int32(0)
	// samplesLeft := int32(totalSamples)

	bpm := 200
	noteLength := 4

	breve := time.Minute / gusic.NoteDuration(bpm) * gusic.NoteDuration(noteLength) * 2
	semibreve := breve / 2
	// minim := semibreve / 2
	// crotchet := semibreve / 4
	quaver := semibreve / 8
	// semiquaver := semibreve / 16
	// demisemiquaver := semibreve / 32

	volume := 0.125

	keys := []Key{}
	whiteKeys := []Key{}
	startOctave := 3
	lastOctave := 5
	octaveCount := lastOctave - startOctave + 1 // +1 because it's inclusive

	// whiteWidth := 70
	whiteWidth := int(width / (int32(octaveCount) * 7)) // 7 is white keys per octave
	blackWidth := int(0.6 * float64(whiteWidth))

	topMargin := 200

	for i := startOctave; i <= lastOctave; i++ {
		// TODO: set duration to 0 and update it based on hold duration
		keys = append(keys,
			NewKey(gusic.C(i, quaver, volume), false),
			NewKey(gusic.CS(i, quaver, volume), true),
			NewKey(gusic.D(i, quaver, volume), false),
			NewKey(gusic.DS(i, quaver, volume), true),
			NewKey(gusic.E(i, quaver, volume), false),
			NewKey(gusic.F(i, quaver, volume), false),
			NewKey(gusic.FS(i, quaver, volume), true),
			NewKey(gusic.G(i, quaver, volume), false),
			NewKey(gusic.GS(i, quaver, volume), true),
			NewKey(gusic.A(i, quaver, volume), false),
			NewKey(gusic.AS(i, quaver, volume), true),
			NewKey(gusic.B(i, quaver, volume), false),
		)
	}

	for _, key := range keys {
		if !key.IsSemitone {
			whiteKeys = append(whiteKeys, key)
		}
	}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Refill audio stream if required
		if rl.IsAudioStreamProcessed(stream) {
			// numSamples := int32(0)
			// if samplesLeft >= maxSamplesPerUpdate {
			// 	numSamples = maxSamplesPerUpdate
			// } else {
			// 	numSamples = samplesLeft
			// }

			// rl.UpdateAudioStream(stream, data[totalSamples-samplesLeft:], numSamples)

			// samplesLeft -= numSamples
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Yellow)
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			// coords := rl.GetMousePosition()
		}

		for i := range whiteKeys {
			rl.DrawRectangle(int32(i*whiteWidth), int32(topMargin), int32(whiteWidth), height-int32(topMargin), rl.White)
			rl.DrawRectangle(int32(i*whiteWidth), int32(topMargin), 2, height-int32(topMargin), rl.Gray)
		}

		for i, key := range keys {
			if key.IsSemitone {
				rl.DrawRectangle(int32(i*blackWidth), int32(topMargin), int32(blackWidth), int32(0.6*float32(height-int32(topMargin))), rl.Black)
			}
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func samplesToFloat32(in []float64) []float32 {
	samples := make([]float32, len(in))
	for i, v := range in {
		samples[i] = float32(v)
	}
	return samples
}
