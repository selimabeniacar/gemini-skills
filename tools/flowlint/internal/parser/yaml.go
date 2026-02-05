package parser

import (
	"gopkg.in/yaml.v3"
)

// Dependencies represents the parsed dependencies.yaml structure
type Dependencies struct {
	Version       string        `yaml:"version"`
	Generated     string        `yaml:"generated"`
	TargetService string        `yaml:"target_service"`
	TargetPath    string        `yaml:"target_path"`
	Documentation Documentation `yaml:"documentation"`
	Service       Service       `yaml:"service"`
	Dependencies  DepsSection   `yaml:"dependencies"`
	Callers       []Caller      `yaml:"callers"`
	External      []External    `yaml:"external"`
	Caches        []Cache       `yaml:"caches"`
}

// Documentation contains references to doc files
type Documentation struct {
	Runbook      string `yaml:"runbook"`
	Architecture string `yaml:"architecture"`
	Readme       string `yaml:"readme"`
	Notes        string `yaml:"notes"`
}

// Service describes the target service
type Service struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Entrypoints []Entrypoint `yaml:"entrypoints"`
}

// Entrypoint describes how traffic enters the service
type Entrypoint struct {
	Type    string   `yaml:"type"`
	Name    string   `yaml:"name"`
	Proto   string   `yaml:"proto"`
	Path    string   `yaml:"path"`
	Topic   string   `yaml:"topic"`
	Methods []string `yaml:"methods"`
}

// DepsSection contains sync and async dependencies
type DepsSection struct {
	Sync  []SyncDep  `yaml:"sync"`
	Async []AsyncDep `yaml:"async"`
}

// SyncDep represents a synchronous dependency (gRPC, HTTP, DB)
type SyncDep struct {
	Name            string   `yaml:"name"`
	Type            string   `yaml:"type"`
	Proto           string   `yaml:"proto"`
	MethodsCalled   []string `yaml:"methods_called"`
	BaseURL         string   `yaml:"base_url"`
	EndpointsCalled []string `yaml:"endpoints_called"`
	Connection      string   `yaml:"connection"`
	Operations      []string `yaml:"operations"`
	Purpose         string   `yaml:"purpose"`
	SourceFile      string   `yaml:"source_file"`
	SourceLine      int      `yaml:"source_line"`
	Timeout         string   `yaml:"timeout"`
	Retry           bool     `yaml:"retry"`
	CircuitBreaker  bool     `yaml:"circuit_breaker"`
}

// AsyncDep represents an asynchronous dependency (Kafka)
type AsyncDep struct {
	Name          string `yaml:"name"`
	Type          string `yaml:"type"`
	Direction     string `yaml:"direction"`
	ConsumerGroup string `yaml:"consumer_group"`
	DLQ           bool   `yaml:"dlq"`
	DLQTopic      string `yaml:"dlq_topic"`
	SourceFile    string `yaml:"source_file"`
	SourceLine    int    `yaml:"source_line"`
}

// Caller represents a service that calls this service
type Caller struct {
	Name          string   `yaml:"name"`
	Type          string   `yaml:"type"`
	MethodsCalled []string `yaml:"methods_called"`
	Source        string   `yaml:"source"`
}

// External represents a third-party external system
type External struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	Purpose    string `yaml:"purpose"`
	SourceFile string `yaml:"source_file"`
	SourceLine int    `yaml:"source_line"`
}

// Cache represents a cache system
type Cache struct {
	Name       string `yaml:"name"`
	Purpose    string `yaml:"purpose"`
	SourceFile string `yaml:"source_file"`
	SourceLine int    `yaml:"source_line"`
}

// ParseDependencies parses YAML content into Dependencies struct
func ParseDependencies(content []byte) (*Dependencies, error) {
	var deps Dependencies
	if err := yaml.Unmarshal(content, &deps); err != nil {
		return nil, err
	}
	return &deps, nil
}
