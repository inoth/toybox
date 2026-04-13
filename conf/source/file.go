package source

import (
	"context"
	"os"
	"time"
)

// FileSource loads config from a local file.
type FileSource struct {
	path string
}

// NewFile creates a file-based config source.
func NewFile(path string) *FileSource {
	return &FileSource{path: path}
}

// Load reads the entire file content as a config string.
func (f *FileSource) Load(_ string) (string, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FileWatcher watches a local file for modifications by polling.
type FileWatcher struct {
	path     string
	interval time.Duration
}

// NewFileWatcher creates a file watcher with an optional polling interval (default 5s).
func NewFileWatcher(path string, interval ...time.Duration) *FileWatcher {
	d := 5 * time.Second
	if len(interval) > 0 {
		d = interval[0]
	}
	return &FileWatcher{path: path, interval: d}
}

// Watche polls the file for modification time changes and sends a trigger on change.
func (w *FileWatcher) Watche(ctx context.Context, trigger chan<- struct{}) {
	var lastMod time.Time
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			info, err := os.Stat(w.path)
			if err != nil {
				continue
			}
			if mod := info.ModTime(); mod.After(lastMod) {
				if !lastMod.IsZero() {
					select {
					case trigger <- struct{}{}:
					default:
					}
				}
				lastMod = mod
			}
		}
	}
}
