package config

import (
	"fmt"
	"io/ioutil"
	"time"

	prometheus_discovery_config "github.com/prometheus/prometheus/discovery/config"
	"github.com/prometheus/prometheus/pkg/relabel"

	yaml "gopkg.in/yaml.v2"
)

type PeriskopConfig struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name             string                                             `yaml:"name"`
	ServiceDiscovery prometheus_discovery_config.ServiceDiscoveryConfig `yaml:",inline"`
	Scraper          Scraper                                            `yaml:"scraper"`
	RelabelConfigs   []*relabel.Config                                  `yaml:"relabel_configs,omitempty"`
}

type Scraper struct {
	RefreshInterval time.Duration `yaml:"refresh_interval"`
	Endpoint        string        `yaml:"endpoint"`
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*PeriskopConfig, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, fmt.Errorf("parsing YAML file %s: %v", filename, err)
	}

	return cfg, nil
}

// Load parses the YAML input s into a Config.
func Load(s string) (*PeriskopConfig, error) {
	cfg := &PeriskopConfig{}

	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
