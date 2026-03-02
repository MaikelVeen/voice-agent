package cesp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	gomp3 "github.com/hajimehoshi/go-mp3"
	"github.com/jfreymuth/oggvorbis"

	"github.com/go-audio/wav"
)

type DecodedAudio struct {
	PCM        io.Reader
	SampleRate int
	Channels   int
}

func DecodeFile(filePath string) (*DecodedAudio, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".wav":
		return decodeWAV(filePath)
	case ".mp3":
		return decodeMP3(filePath)
	case ".ogg":
		return decodeOGG(filePath)
	default:
		return nil, fmt.Errorf("unsupported audio format: %q", ext)
	}
}

func decodeWAV(filePath string) (*DecodedAudio, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening WAV %q: %w", filePath, err)
	}
	defer f.Close()

	dec := wav.NewDecoder(f)
	buf, err := dec.FullPCMBuffer()
	if err != nil {
		return nil, fmt.Errorf("decoding WAV %q: %w", filePath, err)
	}

	pcm := intBufferToInt16LE(buf.Data)

	return &DecodedAudio{
		PCM:        bytes.NewReader(pcm),
		SampleRate: buf.Format.SampleRate,
		Channels:   buf.Format.NumChannels,
	}, nil
}

func intBufferToInt16LE(samples []int) []byte {
	out := make([]byte, len(samples)*2)
	for i, s := range samples {
		// Clamp to int16 range
		if s > 32767 {
			s = 32767
		} else if s < -32768 {
			s = -32768
		}
		binary.LittleEndian.PutUint16(out[i*2:], uint16(int16(s)))
	}
	return out
}

func decodeMP3(filePath string) (*DecodedAudio, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening MP3 %q: %w", filePath, err)
	}
	defer f.Close()

	dec, err := gomp3.NewDecoder(f)
	if err != nil {
		return nil, fmt.Errorf("creating MP3 decoder for %q: %w", filePath, err)
	}

	pcmData, err := io.ReadAll(dec)
	if err != nil {
		return nil, fmt.Errorf("decoding MP3 %q: %w", filePath, err)
	}

	return &DecodedAudio{
		PCM:        bytes.NewReader(pcmData),
		SampleRate: dec.SampleRate(),
		Channels:   2, // go-mp3 always outputs stereo
	}, nil
}

func decodeOGG(filePath string) (*DecodedAudio, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening OGG %q: %w", filePath, err)
	}
	defer f.Close()

	samples, format, err := oggvorbis.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("decoding OGG %q: %w", filePath, err)
	}

	out := make([]byte, len(samples)*2)
	for i, s := range samples {
		// Clamp float32 to [-1, 1] then scale to int16
		if s > 1.0 {
			s = 1.0
		} else if s < -1.0 {
			s = -1.0
		}
		v := int16(s * 32767)
		binary.LittleEndian.PutUint16(out[i*2:], uint16(v))
	}

	return &DecodedAudio{
		PCM:        bytes.NewReader(out),
		SampleRate: format.SampleRate,
		Channels:   format.Channels,
	}, nil
}
