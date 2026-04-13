package bootstrap

import (
	"fmt"
	"log"

	"github.com/inoth/toybox/conf"
	"github.com/inoth/toybox/conf/source"
	"github.com/inoth/toybox/registry"
	"github.com/inoth/toybox/registry/consul"
	"github.com/inoth/toybox/registry/etcd"
	"github.com/inoth/toybox/registry/zookeeper"

	zk "github.com/go-zookeeper/zk"
	consulapi "github.com/hashicorp/consul/api"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// SetupRegistry creates a Registrar from bootstrap config.
// Returns nil if no registry endpoints are configured.
func SetupRegistry(cfg *Config) (registry.Registrar, error) {
	if len(cfg.RegistryEndpoints) == 0 {
		return nil, nil
	}
	switch cfg.RegistryType {
	case "etcd", "":
		return setupEtcd(cfg)
	case "consul":
		return setupConsul(cfg)
	case "zookeeper", "zk":
		return setupZookeeper(cfg)
	default:
		return nil, fmt.Errorf("unsupported registry type: %s", cfg.RegistryType)
	}
}

func setupEtcd(cfg *Config) (registry.Registrar, error) {
	c := clientv3.Config{
		Endpoints:   cfg.RegistryEndpoints,
		DialTimeout: cfg.RegistryTimeout.Duration,
	}
	if cfg.RegistryUsername != "" {
		c.Username = cfg.RegistryUsername
		c.Password = cfg.RegistryPassword
	}
	client, err := clientv3.New(c)
	if err != nil {
		return nil, fmt.Errorf("create etcd client: %w", err)
	}
	return etcd.New(client), nil
}

func setupConsul(cfg *Config) (registry.Registrar, error) {
	c := consulapi.DefaultConfig()
	if len(cfg.RegistryEndpoints) > 0 {
		c.Address = cfg.RegistryEndpoints[0]
	}
	if cfg.RegistryPassword != "" {
		c.Token = cfg.RegistryPassword
	}
	client, err := consulapi.NewClient(c)
	if err != nil {
		return nil, fmt.Errorf("create consul client: %w", err)
	}
	return consul.New(client), nil
}

func setupZookeeper(cfg *Config) (registry.Registrar, error) {
	conn, _, err := zk.Connect(cfg.RegistryEndpoints, cfg.RegistryTimeout.Duration)
	if err != nil {
		return nil, fmt.Errorf("connect zookeeper: %w", err)
	}
	if cfg.RegistryUsername != "" {
		auth := fmt.Sprintf("%s:%s", cfg.RegistryUsername, cfg.RegistryPassword)
		if err := conn.AddAuth("digest", []byte(auth)); err != nil {
			conn.Close()
			return nil, fmt.Errorf("zookeeper auth: %w", err)
		}
	}
	return zookeeper.New(conn), nil
}

// SetupConfig creates a ConfigMate from bootstrap config.
// Priority: remote config endpoint > local config file.
// Returns nil if neither is configured.
func SetupConfig(cfg *Config) (conf.ConfigMate, error) {
	format := conf.Format(cfg.ConfigFormat)
	if format == "" {
		format = conf.FormatYAML
	}

	// Remote config
	if cfg.ConfigEndpoint != "" {
		opts := []source.RemoteOption{}
		if cfg.ConfigToken != "" {
			opts = append(opts, source.WithHeaders(map[string]string{
				"Authorization": "Bearer " + cfg.ConfigToken,
			}))
		}
		remoteSrc := source.NewRemote(cfg.ConfigEndpoint, opts...)
		watcher := source.NewRemoteWatcher(remoteSrc)
		mgr, err := conf.NewManager(remoteSrc, conf.WithFormat(format), conf.WithWatcher(watcher))
		if err != nil {
			return nil, fmt.Errorf("create remote config manager: %w", err)
		}
		log.Printf("Config loaded from remote: %s", cfg.ConfigEndpoint)
		return mgr, nil
	}

	// Local config file
	if cfg.ConfigFile != "" {
		fileSrc := source.NewFile(cfg.ConfigFile)
		fileWatcher := source.NewFileWatcher(cfg.ConfigFile)
		mgr, err := conf.NewManager(fileSrc, conf.WithFormat(format), conf.WithWatcher(fileWatcher))
		if err != nil {
			return nil, fmt.Errorf("create file config manager: %w", err)
		}
		log.Printf("Config loaded from file: %s", cfg.ConfigFile)
		return mgr, nil
	}

	return nil, nil
}
