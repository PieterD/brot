package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ConfigFile string
	OutputFile string
}

func NewConfigFromFlags() (Config, bool) {
	var cfg Config
	flag.StringVar(&cfg.ConfigFile, "config", "config.json", "Filename of the input config file")
	flag.StringVar(&cfg.OutputFile, "output", "mandelbrot.gif", "Filename of the output GIF file")
	flag.Parse()
	if err := cfg.Validate(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "validating config: %v\n", err)
		flag.PrintDefaults()
		return Config{}, false
	}
	return cfg, true
}

func (cfg Config) Validate() error {
	if cfg.ConfigFile == "" {
		return fmt.Errorf("missing -config")
	}
	if cfg.OutputFile == "" {
		return fmt.Errorf("missing -output")
	}
	return nil
}
