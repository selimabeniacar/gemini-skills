package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// DepsFile represents the parsed dependencies.yaml structure
type DepsFile struct {
	Generated string         `yaml:"generated"`
	Services  []ServiceEntry `yaml:"services"`
}

// ServiceEntry represents a single service in the services array
type ServiceEntry struct {
	Name          string         `yaml:"name"`
	TargetPath    string         `yaml:"target_path"`
	Description   string         `yaml:"description"`
	Documentation Documentation  `yaml:"documentation"`
	Entrypoints   []Entrypoint   `yaml:"entrypoints"`
	Dependencies  DepsSection    `yaml:"dependencies"`
	Databases     []Database     `yaml:"databases"`
	Caches        []Cache        `yaml:"caches"`
	External      []External     `yaml:"external"`
	InternalSteps []InternalStep `yaml:"internal_steps"`
}

// Documentation contains references to doc files
type Documentation struct {
	Runbook      string `yaml:"runbook"`
	Architecture string `yaml:"architecture"`
	Notes        string `yaml:"notes"`
}

// Entrypoint describes how traffic enters the service
type Entrypoint struct {
	Type    string   `yaml:"type"`
	Name    string   `yaml:"name"`
	Methods []string `yaml:"methods"`
}

// DepsSection contains sync and async dependencies
type DepsSection struct {
	Sync  []SyncDep  `yaml:"sync"`
	Async []AsyncDep `yaml:"async"`
}

// SyncDep represents a synchronous dependency (gRPC, HTTP)
type SyncDep struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	SourceFile string `yaml:"source_file"`
	SourceLine int    `yaml:"source_line"`
}

// AsyncDep represents an asynchronous dependency (Kafka)
type AsyncDep struct {
	Name       string `yaml:"name"`
	Direction  string `yaml:"direction"`
	SourceFile string `yaml:"source_file"`
	SourceLine int    `yaml:"source_line"`
}

// Database represents a database dependency
type Database struct {
	Name       string `yaml:"name"`
	SourceFile string `yaml:"source_file"`
	SourceLine int    `yaml:"source_line"`
}

// External represents a third-party external system
type External struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
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

// InternalStep represents a processing stage inside a service
type InternalStep struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// ParseDependencies parses YAML content into DepsFile struct.
func ParseDependencies(content []byte) (*DepsFile, error) {
	var deps DepsFile
	if err := yaml.Unmarshal(content, &deps); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if len(deps.Services) == 0 {
		return nil, fmt.Errorf("no services found in dependencies file")
	}

	return &deps, nil
}
