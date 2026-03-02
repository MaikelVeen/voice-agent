package cesp

import (
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
)

const maxSoundFileSize = 1 << 20 // 1 MB

type Pack struct {
	Dir      string
	Manifest *Manifest
}

func LoadPack(nameOrPath string) (*Pack, error) {
	var dir string

	if filepath.IsAbs(nameOrPath) || strings.HasPrefix(nameOrPath, ".") {
		dir = nameOrPath
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("resolving home directory: %w", err)
		}
		dir = filepath.Join(home, ".openpeon", "packs", nameOrPath)
	}

	m, err := LoadManifest(dir)
	if err != nil {
		return nil, err
	}

	if err := m.Validate(); err != nil {
		return nil, fmt.Errorf("invalid manifest: %w", err)
	}

	return &Pack{Dir: dir, Manifest: m}, nil
}

func (p *Pack) ResolveCategory(event string) *Category {
	if cat, ok := p.Manifest.Categories[event]; ok {
		return &cat
	}

	if alias, ok := p.Manifest.CategoryAliases[event]; ok {
		if cat, ok := p.Manifest.Categories[alias]; ok {
			return &cat
		}
	}

	return nil
}

func (p *Pack) PickSound(cat *Category) Sound {
	return cat.Sounds[rand.IntN(len(cat.Sounds))]
}

func (p *Pack) ValidateMagicBytes(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat %q: %w", filePath, err)
	}
	if info.Size() > maxSoundFileSize {
		return fmt.Errorf("file %q exceeds 1 MB limit (%d bytes)", filePath, info.Size())
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening %q: %w", filePath, err)
	}
	defer f.Close()

	var header [4]byte
	if err := binary.Read(f, binary.BigEndian, &header); err != nil {
		return fmt.Errorf("reading header of %q: %w", filePath, err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".wav":
		// RIFF....
		if header[0] != 'R' || header[1] != 'I' || header[2] != 'F' || header[3] != 'F' {
			return fmt.Errorf("%q does not have a valid WAV magic (RIFF)", filePath)
		}
	case ".mp3":
		// ID3 tag header or sync bytes (0xFF 0xFB / 0xFF 0xF3 / 0xFF 0xF2)
		isID3 := header[0] == 'I' && header[1] == 'D' && header[2] == '3'
		isSyncFrame := header[0] == 0xFF && (header[1]&0xE0 == 0xE0)
		if !isID3 && !isSyncFrame {
			return fmt.Errorf("%q does not have a valid MP3 magic (ID3 or sync frame)", filePath)
		}
	case ".ogg":
		// OggS
		if header[0] != 'O' || header[1] != 'g' || header[2] != 'g' || header[3] != 'S' {
			return fmt.Errorf("%q does not have a valid OGG magic (OggS)", filePath)
		}
	default:
		return fmt.Errorf("unsupported audio extension %q", ext)
	}

	return nil
}

func (p *Pack) SoundPath(s Sound) string {
	return filepath.Join(p.Dir, filepath.FromSlash(s.File))
}
