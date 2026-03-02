package audio

import (
	"fmt"
	"io"
	"time"

	"github.com/ebitengine/oto/v3"
)

func PlayPCM(pcm io.Reader) error {
	otoCtx, readyChan, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   24000,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		return fmt.Errorf("initializing audio context: %w", err)
	}
	<-readyChan

	player := otoCtx.NewPlayer(pcm)
	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// IsPlaying() returns false when oto's software buffer is empty, but the
	// hardware driver buffer (AudioQueue on macOS: 4x12288 bytes) may still
	// hold ~256ms of audio at 24kHz mono int16. Sleep to allow hardware drain.
	// See: https://github.com/ebitengine/oto/issues/237
	time.Sleep(300 * time.Millisecond)

	return nil
}

func PlayPCMWithFormat(pcm io.Reader, sampleRate, channels int) error {
	otoCtx, readyChan, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: channels,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		return fmt.Errorf("initializing audio context: %w", err)
	}
	<-readyChan

	player := otoCtx.NewPlayer(pcm)
	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// Same hardware drain buffer concern as PlayPCM — see comment above.
	time.Sleep(300 * time.Millisecond)

	return nil
}
