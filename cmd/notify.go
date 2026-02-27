package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/MaikelVeen/voice-agent/internal/tts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type notifyPayload struct {
	SessionID        string `json:"session_id"`
	TranscriptPath   string `json:"transcript_path"`
	Cwd              string `json:"cwd"`
	PermissionMode   string `json:"permission_mode"`
	HookEventName    string `json:"hook_event_name"`
	Message          string `json:"message"`
	Title            string `json:"title"`
	NotificationType string `json:"notification_type"`
}

type NotifyCommand struct {
	Command *cobra.Command
}

func NewNotifyCommand() *NotifyCommand {
	cmd := &cobra.Command{
		Use:   "notify",
		Short: "Speak a Claude Code notification message",
		Long:  "Reads a Claude Code Notification hook payload from stdin and speaks the message field using OpenAI TTS.",
		RunE:  runNotify,
	}

	flags := cmd.Flags()
	flags.String(argSpeakVoice, "alloy", "Voice to use (alloy, ash, coral, echo, sage, shimmer, verse)")
	_ = viper.BindPFlag(argSpeakVoice, flags.Lookup(argSpeakVoice))

	return &NotifyCommand{Command: cmd}
}

func runNotify(_ *cobra.Command, _ []string) error {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	var payload notifyPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("parsing notification payload: %w", err)
	}

	if payload.Message == "" {
		return fmt.Errorf("notification payload contains no message")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	voice := viper.GetString(argSpeakVoice)

	return tts.Speak(apiKey, payload.Message, voice)
}
