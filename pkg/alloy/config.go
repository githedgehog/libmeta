// Copyright 2025 Hedgehog
// SPDX-License-Identifier: Apache-2.0

package alloy

import (
	_ "embed"
	"fmt"
	"strings"

	"go.githedgehog.com/libmeta/pkg/tmpl"
)

type Config struct {
	Hostname string
	Targets  Targets
	ProxyURL string

	Scrapes  map[string]Scrape
	LogFiles []string
	Kube     Kube
}

type Scrape struct {
	IntervalSeconds uint

	Address string
	Self    ScrapeSelf
	Unix    ScrapeUnix
}

type ScrapeSelf struct {
	Enable bool
}

type ScrapeUnix struct {
	Enable     bool
	Collectors []string
}

type Kube struct {
	PodLogs bool
	Events  bool
}

func (cfg *Config) Validate() error {
	if cfg == nil {
		return fmt.Errorf("config is nil") //nolint:err113
	}

	if err := cfg.Targets.Validate(); err != nil {
		return fmt.Errorf("invalid targets: %w", err)
	}

	for name, scrape := range cfg.Scrapes {
		if err := validateIdentifier(name); err != nil {
			return fmt.Errorf("invalid scrape name %q: %w", name, err)
		}

		if err := scrape.Validate(); err != nil {
			return fmt.Errorf("invalid scrape %q: %w", name, err)
		}
	}

	return nil
}

func (s *Scrape) Validate() error {
	if s == nil {
		return fmt.Errorf("scrape is nil") //nolint:err113
	}

	opts := 0
	if s.Address != "" {
		opts++
	}
	if s.Self.Enable {
		opts++
	}
	if s.Unix.Enable {
		opts++
	}
	if opts == 0 {
		return fmt.Errorf("no scrape options enabled") //nolint:err113
	}
	if opts > 1 {
		return fmt.Errorf("multiple scrape options enabled") //nolint:err113
	}

	return nil
}

func (k *Kube) Validate() error {
	if k == nil {
		return fmt.Errorf("kube is nil") //nolint:err113
	}

	return nil
}

//go:embed config.alloy.tmpl
var configTemplate string

func (cfg *Config) Render() ([]byte, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil") //nolint:err113
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	data, err := tmpl.Render("config.alloy.tmpl", configTemplate, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to render config: %w", err)
	}

	var res strings.Builder
	for line := range strings.Lines(string(data)) {
		if strings.TrimSpace(line) == "" {
			continue
		}
		res.WriteString(line)
	}

	return []byte(res.String()), nil
}
