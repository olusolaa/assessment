package config

import (
	"encoding/json"
	"os"
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
	
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	
	decoder := json.NewDecoder(file)
	return decoder.Decode(c)
}


func (c *Config) SaveToFile(filename string) error {
	
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
} 
