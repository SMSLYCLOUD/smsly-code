package cicd

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type WorkflowDefinition struct {
	Name string                   `yaml:"name"`
	On   interface{}              `yaml:"on"`
	Jobs map[string]JobDefinition `yaml:"jobs"`
}

type JobDefinition struct {
	Name   string           `yaml:"name,omitempty"`
	RunsOn string           `yaml:"runs-on"`
	Steps  []StepDefinition `yaml:"steps"`
	Needs  interface{}      `yaml:"needs,omitempty"` // string or []string
}

type StepDefinition struct {
	Name string            `yaml:"name,omitempty"`
	Uses string            `yaml:"uses,omitempty"`
	Run  string            `yaml:"run,omitempty"`
	Env  map[string]string `yaml:"env,omitempty"`
}

// Parse parses the workflow YAML content.
func Parse(content []byte) (*WorkflowDefinition, error) {
	var wf WorkflowDefinition
	if err := yaml.Unmarshal(content, &wf); err != nil {
		return nil, fmt.Errorf("failed to parse workflow yaml: %w", err)
	}

	if len(wf.Jobs) == 0 {
		return nil, fmt.Errorf("workflow must have at least one job")
	}

	return &wf, nil
}
