package bootstrap

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Config holds bootstrap-level configuration that must be resolved
// before connecting to remote config centers or service registries.
type Config struct {
	// Service identity
	ServiceName    string `json:"service_name" yaml:"service_name" toml:"service_name"`
	ServiceVersion string `json:"service_version" yaml:"service_version" toml:"service_version"`

	// Registry
	RegistryType      string   `json:"registry_type" yaml:"registry_type" toml:"registry_type"`                // etcd, consul, zookeeper
	RegistryEndpoints []string `json:"registry_endpoints" yaml:"registry_endpoints" toml:"registry_endpoints"` // e.g. ["localhost:2379"]
	RegistryTimeout   Duration `json:"registry_timeout" yaml:"registry_timeout" toml:"registry_timeout"`
	RegistryUsername  string   `json:"registry_username" yaml:"registry_username" toml:"registry_username"`
	RegistryPassword  string   `json:"registry_password" yaml:"registry_password" toml:"registry_password"`

	// Remote config source
	ConfigEndpoint string `json:"config_endpoint" yaml:"config_endpoint" toml:"config_endpoint"` // remote config center URL
	ConfigFormat   string `json:"config_format" yaml:"config_format" toml:"config_format"`       // yaml, json, toml
	ConfigToken    string `json:"config_token" yaml:"config_token" toml:"config_token"`

	// Local config fallback
	ConfigFile string `json:"config_file" yaml:"config_file" toml:"config_file"` // local config file path
}

// Duration wraps time.Duration for text unmarshaling.
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.Duration.String()), nil
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		RegistryType:    "etcd",
		RegistryTimeout: Duration{5 * time.Second},
		ConfigFormat:    "yaml",
	}
}

// Load resolves bootstrap config by merging (in priority order):
//  1. Command-line flags (highest)
//  2. Environment variables
//  3. Defaults (lowest)
func Load() *Config {
	cfg := Default()
	cfg.loadEnv()
	cfg.loadFlags()
	return cfg
}

// loadEnv reads configuration from environment variables.
// Env var naming: TOYBOX_<UPPER_SNAKE_FIELD>
func (c *Config) loadEnv() {
	if v := os.Getenv("TOYBOX_SERVICE_NAME"); v != "" {
		c.ServiceName = v
	}
	if v := os.Getenv("TOYBOX_SERVICE_VERSION"); v != "" {
		c.ServiceVersion = v
	}
	if v := os.Getenv("TOYBOX_REGISTRY_TYPE"); v != "" {
		c.RegistryType = v
	}
	if v := os.Getenv("TOYBOX_REGISTRY_ENDPOINTS"); v != "" {
		c.RegistryEndpoints = splitComma(v)
	}
	if v := os.Getenv("TOYBOX_REGISTRY_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			c.RegistryTimeout = Duration{d}
		}
	}
	if v := os.Getenv("TOYBOX_REGISTRY_USERNAME"); v != "" {
		c.RegistryUsername = v
	}
	if v := os.Getenv("TOYBOX_REGISTRY_PASSWORD"); v != "" {
		c.RegistryPassword = v
	}
	if v := os.Getenv("TOYBOX_CONFIG_ENDPOINT"); v != "" {
		c.ConfigEndpoint = v
	}
	if v := os.Getenv("TOYBOX_CONFIG_FORMAT"); v != "" {
		c.ConfigFormat = v
	}
	if v := os.Getenv("TOYBOX_CONFIG_TOKEN"); v != "" {
		c.ConfigToken = v
	}
	if v := os.Getenv("TOYBOX_CONFIG_FILE"); v != "" {
		c.ConfigFile = v
	}
}

// loadFlags reads configuration from command-line flags.
// Flags take precedence over env vars.
func (c *Config) loadFlags() {
	fs := flag.NewFlagSet("toybox", flag.ContinueOnError)

	serviceName := fs.String("service-name", "", "Service name")
	serviceVersion := fs.String("service-version", "", "Service version")
	registryType := fs.String("registry-type", "", "Registry type: etcd, consul, zookeeper")
	registryEndpoints := fs.String("registry-endpoints", "", "Registry endpoints (comma-separated)")
	registryTimeout := fs.String("registry-timeout", "", "Registry dial timeout (e.g. 5s)")
	registryUsername := fs.String("registry-username", "", "Registry auth username")
	registryPassword := fs.String("registry-password", "", "Registry auth password")
	configEndpoint := fs.String("config-endpoint", "", "Remote config center URL")
	configFormat := fs.String("config-format", "", "Config format: yaml, json, toml")
	configToken := fs.String("config-token", "", "Remote config auth token")
	configFile := fs.String("config-file", "", "Local config file path")

	// Silently ignore unknown flags (allows coexistence with app flags).
	fs.SetOutput(nil)
	_ = fs.Parse(os.Args[1:])

	if *serviceName != "" {
		c.ServiceName = *serviceName
	}
	if *serviceVersion != "" {
		c.ServiceVersion = *serviceVersion
	}
	if *registryType != "" {
		c.RegistryType = *registryType
	}
	if *registryEndpoints != "" {
		c.RegistryEndpoints = splitComma(*registryEndpoints)
	}
	if *registryTimeout != "" {
		if d, err := time.ParseDuration(*registryTimeout); err == nil {
			c.RegistryTimeout = Duration{d}
		}
	}
	if *registryUsername != "" {
		c.RegistryUsername = *registryUsername
	}
	if *registryPassword != "" {
		c.RegistryPassword = *registryPassword
	}
	if *configEndpoint != "" {
		c.ConfigEndpoint = *configEndpoint
	}
	if *configFormat != "" {
		c.ConfigFormat = *configFormat
	}
	if *configToken != "" {
		c.ConfigToken = *configToken
	}
	if *configFile != "" {
		c.ConfigFile = *configFile
	}
}

// String returns a redacted summary for logging.
func (c *Config) String() string {
	endpoints := strings.Join(c.RegistryEndpoints, ",")
	return fmt.Sprintf(
		"service=%s/%s registry=%s(%s) config_endpoint=%s config_file=%s format=%s",
		c.ServiceName, c.ServiceVersion,
		c.RegistryType, endpoints,
		c.ConfigEndpoint, c.ConfigFile, c.ConfigFormat,
	)
}

func splitComma(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
