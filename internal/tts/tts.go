package tts

import (
	"context"
	"fmt"

	"github.com/MaikelVeen/voice-agent/internal/audio"
	"github.com/MaikelVeen/voice-agent/internal/logging"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func Speak(apiKey, text, voice string, entry *logging.Entry) error {
	client := openai.NewClient(option.WithAPIKey(apiKey))

	res, err := client.Audio.Speech.New(context.Background(), openai.AudioSpeechNewParams{
		Model:          openai.SpeechModelGPT4oMiniTTS,
		Input:          text,
		Voice:          openai.AudioSpeechNewParamsVoice(voice),
		ResponseFormat: openai.AudioSpeechNewParamsResponseFormatPCM,
	})
	if err != nil {
		return fmt.Errorf("requesting speech: %w", err)
	}
	defer res.Body.Close()

	entry.CaptureHeaders(res.Header)
	if err := logging.Write(entry); err != nil {
		// Log failures are non-fatal
		fmt.Printf("warning: failed to write log entry: %v\n", err)
	}

	if err := audio.PlayPCM(res.Body); err != nil {
		return fmt.Errorf("playing audio: %w", err)
	}

	return nil
}
