package cmd

import (
	"fmt"

	"github.com/MaikelVeen/voice-agent/internal/audio"
	"github.com/MaikelVeen/voice-agent/internal/cesp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	argPlayPack  = "pack"
	argPlayEvent = "event"
)

type PlayCommand struct {
	Command *cobra.Command
}

func NewPlayCommand() *PlayCommand {
	cmd := &cobra.Command{
		Use:   "play",
		Short: "Play a sound from a CESP sound pack",
		Long:  "Loads a CESP v1.0 sound pack and plays a sound for the given event category.",
		RunE:  runPlay,
	}

	flags := cmd.Flags()
	flags.String(argPlayPack, "", "Pack name (resolved from ~/.openpeon/packs/) or path to pack directory")
	_ = viper.BindPFlag(argPlayPack, flags.Lookup(argPlayPack))
	_ = cmd.MarkFlagRequired(argPlayPack)

	flags.String(argPlayEvent, "", "CESP event category (e.g. task.complete, task.error)")
	_ = viper.BindPFlag(argPlayEvent, flags.Lookup(argPlayEvent))
	_ = cmd.MarkFlagRequired(argPlayEvent)

	return &PlayCommand{Command: cmd}
}

func runPlay(_ *cobra.Command, _ []string) error {
	packArg := viper.GetString(argPlayPack)
	event := viper.GetString(argPlayEvent)

	pack, err := cesp.LoadPack(packArg)
	if err != nil {
		return fmt.Errorf("loading pack %q: %w", packArg, err)
	}

	cat := pack.ResolveCategory(event)
	if cat == nil {
		// Per spec: missing categories are silent
		return nil
	}

	sound := pack.PickSound(cat)
	soundPath := pack.SoundPath(sound)

	if err := pack.ValidateMagicBytes(soundPath); err != nil {
		return fmt.Errorf("validating sound file: %w", err)
	}

	decoded, err := cesp.DecodeFile(soundPath)
	if err != nil {
		return fmt.Errorf("decoding sound file: %w", err)
	}

	return audio.PlayPCMWithFormat(decoded.PCM, decoded.SampleRate, decoded.Channels)
}
