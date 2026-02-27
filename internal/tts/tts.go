package tts

import (
	"context"
	"fmt"

	"github.com/MaikelVeen/voice-agent/internal/audio"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func Speak(apiKey, text, voice string) error {
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

	if err := audio.PlayPCM(res.Body); err != nil {
		return fmt.Errorf("playing audio: %w", err)
	}

	return nil
}
