package cicd

import (
	"testing"
)

func TestParse_Valid(t *testing.T) {
	yamlContent := `
name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Test
        run: go test ./...
`
	wf, err := Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if wf.Name != "CI" {
		t.Errorf("Expected name CI, got %s", wf.Name)
	}
	if len(wf.Jobs) != 1 {
		t.Errorf("Expected 1 job, got %d", len(wf.Jobs))
	}
	job, ok := wf.Jobs["build"]
	if !ok {
		t.Fatalf("Job 'build' not found")
	}
	if job.RunsOn != "ubuntu-latest" {
		t.Errorf("Expected runs-on ubuntu-latest, got %s", job.RunsOn)
	}
	if len(job.Steps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(job.Steps))
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	yamlContent := `
name: CI
on: [push
`
	_, err := Parse([]byte(yamlContent))
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestParse_NoJobs(t *testing.T) {
	yamlContent := `
name: CI
on: push
`
	_, err := Parse([]byte(yamlContent))
	if err == nil {
		t.Error("Expected error for no jobs, got nil")
	} else if err.Error() != "workflow must have at least one job" {
		t.Errorf("Expected 'workflow must have at least one job', got %v", err)
	}
}
