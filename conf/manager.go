package conf

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"sync"
)

// ManagerOption configures a Manager.
type ManagerOption func(*Manager)

// WithFormat sets the config format (yaml, json, toml). Default is yaml.
func WithFormat(format Format) ManagerOption {
	return func(m *Manager) {
		m.format = format
	}
}

// WithWatcher sets a Watcher for config hot-reload.
func WithWatcher(w Watcher) ManagerOption {
	return func(m *Manager) {
		m.watcher = w
	}
}

// Manager implements ConfigMate and provides config loading, decoding and hot-reload.
type Manager struct {
	mu       sync.RWMutex
	source   Source
	watcher  Watcher
	format   Format
	codec    Codec
	raw      []byte
	sections map[string]any
	onChange []func()
}

// NewManager creates a Manager that loads config from the given Source.
func NewManager(source Source, opts ...ManagerOption) (*Manager, error) {
	m := &Manager{
		source: source,
		format: FormatYAML,
	}
	for _, opt := range opts {
		opt(m)
	}
	codec, err := NewCodec(m.format)
	if err != nil {
		return nil, err
	}
	m.codec = codec
	if err := m.load(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) load() error {
	data, err := m.source.Load("*")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.raw = []byte(data)
	m.sections = make(map[string]any)
	return m.codec.Unmarshal(m.raw, &m.sections)
}

// PrimitiveDecode decodes config sections into matching ConfigureMatcher instances.
func (m *Manager) PrimitiveDecode(vals ...ConfigureMatcher) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range vals {
		name := v.TransportName()
		raw, ok := m.sections[name]
		if !ok {
			continue
		}
		// Re-encode the section then decode into the target struct
		buf, err := m.codec.Marshal(raw)
		if err != nil {
			return fmt.Errorf("re-encode section %q: %w", name, err)
		}
		if err := m.codec.Unmarshal(buf, v); err != nil {
			return fmt.Errorf("decode section %q: %w", name, err)
		}
	}
	return nil
}

// Get returns the raw value for a config section.
func (m *Manager) Get(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.sections[key]
	return v, ok
}

// Decode unmarshals a config section into the given target.
func (m *Manager) Decode(key string, target interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	raw, ok := m.sections[key]
	if !ok {
		return fmt.Errorf("config section %q not found", key)
	}
	buf, err := m.codec.Marshal(raw)
	if err != nil {
		return fmt.Errorf("re-encode section %q: %w", key, err)
	}
	return m.codec.Unmarshal(buf, target)
}

// Raw returns the full raw config bytes.
func (m *Manager) Raw() []byte {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.raw
}

// Checksum returns the SHA-256 checksum of the current config.
func (m *Manager) Checksum() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h := sha256.Sum256(m.raw)
	return fmt.Sprintf("%x", h)
}

// OnChange registers a callback that fires when config is reloaded.
func (m *Manager) OnChange(fn func()) {
	m.onChange = append(m.onChange, fn)
}

// Watch starts watching for config changes and reloads automatically.
func (m *Manager) Watch(ctx context.Context) error {
	if m.watcher == nil {
		return nil
	}
	trigger := make(chan struct{}, 1)
	go m.watcher.Watche(ctx, trigger)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-trigger:
			if err := m.load(); err != nil {
				log.Printf("config reload error: %v", err)
				continue
			}
			log.Println("config reloaded successfully")
			for _, fn := range m.onChange {
				fn()
			}
		}
	}
}
