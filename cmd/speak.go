package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MaikelVeen/voice-agent/internal/tts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	argSpeakText  = "text"
	argSpeakVoice = "voice"
)

type SpeakCommand struct {
	Command *cobra.Command
}

func NewSpeakCommand() *SpeakCommand {
	cmd := &cobra.Command{
		Use:   "speak [text...]",
		Short: "Speak text using OpenAI TTS",
		Long:  "Speak text using OpenAI's text-to-speech API. Text can be provided as positional arguments or via the --text flag.",
		RunE:  runSpeak,
	}

	flags := cmd.Flags()
	flags.StringP(argSpeakText, "t", "", "Text to speak (alternative to positional args)")
	_ = viper.BindPFlag(argSpeakText, flags.Lookup(argSpeakText))

	flags.String(argSpeakVoice, "alloy", "Voice to use (alloy, ash, coral, echo, sage, shimmer, verse)")
	_ = viper.BindPFlag(argSpeakVoice, flags.Lookup(argSpeakVoice))

	return &SpeakCommand{Command: cmd}
}

func runSpeak(_ *cobra.Command, args []string) error {
	text := viper.GetString(argSpeakText)
	if text == "" && len(args) > 0 {
		text = strings.Join(args, " ")
	}
	if text == "" {
		return fmt.Errorf("no text provided; pass text as arguments or use --text")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	voice := viper.GetString(argSpeakVoice)

	return tts.Speak(apiKey, text, voice)
}
