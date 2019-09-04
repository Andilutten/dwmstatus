package main

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type (
	// A Config object is read from
	// $HOME/.config/dwmstatus/config.yaml
	Config struct {
		Items []Item `yaml:"items"`
	}

	// An Item describes one single item
	// in the statusbar.
	Item struct {
		Command  string        `yaml:"command"`
		Interval time.Duration `yaml:"interval"`
		Name     string        `yaml:"name"`
	}
)

func NewConfig(path string) (*Config, error) {
	// Read contents from file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Unmarshal contents into
	// config object
	cfg := new(Config)
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
