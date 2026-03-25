package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		MetricsPort int    `yaml:"metrics_port"`
		HTTPPort    int    `yaml:"http_port"`
	} `yaml:"server"`
	Storage struct {
		CleanupInterval  int    `yaml:"cleanup_interval_seconds"`
		MaxMemoryMB      int    `yaml:"max_memory_mb"`
		SnapshotFile     string `yaml:"snapshot_file"`
		SnapshotInterval int    `yaml:"snapshot_interval_seconds"`
	} `yaml:"storage"`
}

func LoadConfig(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	return c, err
}
