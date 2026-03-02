package cesp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	version "github.com/hashicorp/go-version"
)

type Sound struct {
	File   string `json:"file"`
	Label  string `json:"label"`
	SHA256 string `json:"sha256,omitempty"`
	Icon   string `json:"icon,omitempty"`
}

type Category struct {
	Sounds []Sound `json:"sounds"`
	Icon   string  `json:"icon,omitempty"`
}

type Manifest struct {
	CESPVersion     string              `json:"cesp_version"`
	Name            string              `json:"name"`
	DisplayName     string              `json:"display_name"`
	Version         string              `json:"version"`
	Categories      map[string]Category `json:"categories"`
	CategoryAliases map[string]string   `json:"category_aliases,omitempty"`
	Description     string              `json:"description,omitempty"`
}

func LoadManifest(dir string) (*Manifest, error) {
	data, err := os.ReadFile(filepath.Join(dir, "openpeon.json"))
	if err != nil {
		return nil, fmt.Errorf("reading openpeon.json: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing openpeon.json: %w", err)
	}

	return &m, nil
}

func (m *Manifest) Validate() error {
	if m.CESPVersion != "1.0" {
		return fmt.Errorf("unsupported cesp_version %q; expected \"1.0\"", m.CESPVersion)
	}
	if m.Name == "" {
		return fmt.Errorf("manifest missing required field: name")
	}
	if m.DisplayName == "" {
		return fmt.Errorf("manifest missing required field: display_name")
	}
	if m.Version == "" {
		return fmt.Errorf("manifest missing required field: version")
	}
	if _, err := version.NewSemver(m.Version); err != nil {
		return fmt.Errorf("manifest version %q is not valid semver: %w", m.Version, err)
	}
	if len(m.Categories) == 0 {
		return fmt.Errorf("manifest must define at least one category")
	}

	for catName, cat := range m.Categories {
		for i, s := range cat.Sounds {
			if s.File == "" {
				return fmt.Errorf("category %q sound[%d] missing required field: file", catName, i)
			}
			if s.Label == "" {
				return fmt.Errorf("category %q sound[%d] missing required field: label", catName, i)
			}
			if strings.Contains(s.File, "..") {
				return fmt.Errorf("category %q sound[%d] file path contains illegal traversal: %q", catName, i, s.File)
			}
		}
	}

	return nil
}
