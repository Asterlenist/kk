package history

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"
)

type entry struct {
	Timestamp string `json:"timestamp"`
	Command   string `json:"command"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
}

func Append(path string, maxSize int64, command string, commandErr error) error {
	record := entry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Command:   command,
		Status:    "success",
	}
	if commandErr != nil {
		record.Status = "error"
		record.Error = commandErr.Error()
	}

	line, err := encode(record)
	if err != nil {
		return err
	}
	if maxSize > 0 && int64(len(line)) > maxSize {
		record.Error = ""
		record.Command = truncateToSize(record, command, maxSize)
		line, err = encode(record)
		if err != nil {
			return err
		}
		if int64(len(line)) > maxSize {
			return fmt.Errorf("logging.max_size is too small for a history entry")
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create history directory: %w", err)
	}
	if maxSize == 0 {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("open history file %q: %w", path, err)
		}
		if _, err := file.Write(line); err != nil {
			file.Close()
			return fmt.Errorf("append history file %q: %w", path, err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("close history file %q: %w", path, err)
		}
		return nil
	}

	current, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read history file %q: %w", path, err)
	}
	combined := append(current, line...)
	if int64(len(combined)) > maxSize {
		start := int64(len(combined)) - maxSize
		if start > 0 && combined[start-1] != '\n' {
			if newline := bytes.IndexByte(combined[start:], '\n'); newline >= 0 {
				start += int64(newline + 1)
			}
		}
		combined = combined[start:]
	}
	if err := os.WriteFile(path, combined, 0600); err != nil {
		return fmt.Errorf("write history file %q: %w", path, err)
	}
	return nil
}

func encode(record entry) ([]byte, error) {
	line, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("encode history entry: %w", err)
	}
	return append(line, '\n'), nil
}

func truncateToSize(record entry, command string, maxSize int64) string {
	low, high := 0, utf8.RuneCountInString(command)
	for low < high {
		mid := (low + high + 1) / 2
		record.Command = string([]rune(command)[:mid]) + "..."
		line, _ := encode(record)
		if int64(len(line)) <= maxSize {
			low = mid
		} else {
			high = mid - 1
		}
	}
	if low == utf8.RuneCountInString(command) {
		return command
	}
	if low == 0 {
		return ""
	}
	return string([]rune(command)[:low]) + "..."
}
