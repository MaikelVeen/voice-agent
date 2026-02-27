package logging

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	Timestamp        time.Time `json:"timestamp"`
	SessionID        string    `json:"session_id,omitempty"`
	Cwd              string    `json:"cwd,omitempty"`
	HookEventName    string    `json:"hook_event_name,omitempty"`
	NotificationType string    `json:"notification_type,omitempty"`
	Text             string    `json:"text"`
	Voice            string    `json:"voice"`
	Model            string    `json:"model"`
	// Response headers from the OpenAI API call
	RequestID    string `json:"request_id,omitempty"`
	ProcessingMS string `json:"processing_ms,omitempty"`
}

func (e *Entry) CaptureHeaders(h http.Header) {
	e.RequestID = h.Get("X-Request-Id")
	e.ProcessingMS = h.Get("Openai-Processing-Ms")
}

func Write(e *Entry) error {
	e.Timestamp = time.Now().UTC()

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolving home directory: %w", err)
	}

	dir := filepath.Join(home, ".voice-agent")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating log directory: %w", err)
	}

	f, err := os.OpenFile(filepath.Join(dir, "tts.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}
	defer f.Close()

	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("marshaling log entry: %w", err)
	}

	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}
