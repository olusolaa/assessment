package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)


type Config struct {
	
	DisabledSorters []string `json:"disabled_sorters"`
	
	
	DefaultPageSize int `json:"default_page_size"`
}


func NewConfig() *Config {
	return &Config{
		DisabledSorters: []string{},
		DefaultPageSize: 10,
	}
}


func (c *Config) LoadFromFile(filename string) error {
	
	// Validate filename to prevent path traversal
	cleanPath := filepath.Clean(filename)
	if filepath.IsAbs(cleanPath) || strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid filename path: potential directory traversal attempt")
	}
	
	file, err := os.Open(cleanPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	
	decoder := json.NewDecoder(file)
	return decoder.Decode(c)
}


func (c *Config) SaveToFile(filename string) error {
	
	// Validate filename to prevent path traversal
	cleanPath := filepath.Clean(filename)
	if filepath.IsAbs(cleanPath) || strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid filename path: potential directory traversal attempt")
	}
	
	file, err := os.Create(cleanPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
} 
