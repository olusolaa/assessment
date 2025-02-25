package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)


type SorterConfig struct {
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	Ascending bool   `json:"ascending"`
}


type Config struct {
	Sorters map[string]SorterConfig `json:"sorters"`
	mu      sync.RWMutex
}


func NewConfig() *Config {
	return &Config{
		Sorters: map[string]SorterConfig{
			"price": {
				Name:      "Price",
				Enabled:   true,
				Ascending: true,
			},
			"sales_per_view": {
				Name:      "Sales per View",
				Enabled:   true,
				Ascending: false,
			},
			"creation_date": {
				Name:      "Creation Date",
				Enabled:   true,
				Ascending: true,
			},
			"name": {
				Name:      "Name",
				Enabled:   true,
				Ascending: true,
			},
		},
	}
}


func (c *Config) LoadFromFile(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Validate filename to prevent path traversal
	cleanPath := filepath.Clean(filename)
	if filepath.IsAbs(cleanPath) || strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid filename path: potential directory traversal attempt")
	}

	data, err := ioutil.ReadFile(cleanPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	return nil
}


func (c *Config) SaveToFile(filename string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing config: %w", err)
	}

	// Validate filename to prevent path traversal
	cleanPath := filepath.Clean(filename)
	if filepath.IsAbs(cleanPath) || strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid filename path: potential directory traversal attempt")
	}

	// Use more restrictive file permissions (0600 instead of 0644)
	err = ioutil.WriteFile(cleanPath, data, 0600)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}


func (c *Config) GetSorterConfig(key string) (SorterConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	config, exists := c.Sorters[key]
	return config, exists
}


func (c *Config) SetSorterConfig(key string, config SorterConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Sorters[key] = config
}


func (c *Config) IsSorterEnabled(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	config, exists := c.Sorters[key]
	return exists && config.Enabled
}


func (c *Config) GetEnabledSorters() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var enabled []string
	for key, config := range c.Sorters {
		if config.Enabled {
			enabled = append(enabled, key)
		}
	}
	return enabled
} 
