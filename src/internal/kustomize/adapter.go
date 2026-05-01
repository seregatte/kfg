package kustomize

import (
	"bytes"
	"fmt"
	"strings"

	yamlv3 "gopkg.in/yaml.v3"
	"github.com/seregatte/kfg/src/internal/logger"
	"github.com/seregatte/kfg/src/internal/manifest"
	"sigs.k8s.io/kustomize/api/resmap"
)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) ResMapToResources(resMap resmap.ResMap) ([]manifest.ParsedResource, error) {
	if resMap == nil {
		return nil, nil
	}

	yamlBytes, err := resMap.AsYaml()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize ResMap to YAML: %w", err)
	}

	return a.parseMultiDocumentYAML(yamlBytes)
}

func (a *Adapter) parseMultiDocumentYAML(yamlBytes []byte) ([]manifest.ParsedResource, error) {
	var resources []manifest.ParsedResource

	decoder := yamlv3.NewDecoder(bytes.NewReader(yamlBytes))

	for {
		var node yamlv3.Node
		if err := decoder.Decode(&node); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to decode YAML node: %w", err)
		}

		if node.Kind == 0 {
			continue
		}

		res, err := a.parseYamlNode(&node)
		if err != nil {
			return nil, err
		}

		if res.Kind() != "" {
			resources = append(resources, res)
		}
	}

	return resources, nil
}

func (a *Adapter) parseYamlNode(node *yamlv3.Node) (manifest.ParsedResource, error) {
	if node.Kind != yamlv3.DocumentNode {
		return manifest.ParsedResource{}, fmt.Errorf("expected DocumentNode, got %v", node.Kind)
	}

	if len(node.Content) == 0 {
		return manifest.ParsedResource{}, nil
	}

	contentNode := node.Content[0]
	if contentNode.Kind != yamlv3.MappingNode {
		return manifest.ParsedResource{}, fmt.Errorf("expected MappingNode, got %v", contentNode.Kind)
	}

	var kindOnly struct {
		Kind string `yaml:"kind"`
	}
	if err := contentNode.Decode(&kindOnly); err != nil {
		return manifest.ParsedResource{}, fmt.Errorf("failed to decode kind: %w", err)
	}

	switch kindOnly.Kind {
	case "Step":
		var step manifest.Step
		if err := contentNode.Decode(&step); err != nil {
			return manifest.ParsedResource{}, fmt.Errorf("failed to decode Step: %w", err)
		}
		if err := step.Validate(); err != nil {
			return manifest.ParsedResource{}, fmt.Errorf("Step validation failed for %s: %v", step.Metadata.Name, err)
		}
		return manifest.ParsedResource{Step: &step}, nil

	case "Cmd":
		var cmd manifest.Cmd
		if err := contentNode.Decode(&cmd); err != nil {
			return manifest.ParsedResource{}, fmt.Errorf("failed to decode Cmd: %w", err)
		}
		if err := cmd.Validate(); err != nil {
			return manifest.ParsedResource{}, fmt.Errorf("Cmd validation failed for %s: %v", cmd.Metadata.Name, err)
		}
		return manifest.ParsedResource{Cmd: &cmd}, nil

	case "CmdWorkflow":
		var workflow manifest.CmdWorkflow
		if err := contentNode.Decode(&workflow); err != nil {
			return manifest.ParsedResource{}, fmt.Errorf("failed to decode CmdWorkflow: %w", err)
		}
		if err := workflow.Validate(); err != nil {
			return manifest.ParsedResource{}, fmt.Errorf("CmdWorkflow validation failed for %s: %v", workflow.Metadata.Name, err)
		}
		return manifest.ParsedResource{CmdWorkflow: &workflow}, nil

	default:
		return manifest.ParsedResource{}, fmt.Errorf("unsupported kind: %s (supported: %s)", kindOnly.Kind, strings.Join(manifest.SupportedKinds, ", "))
	}
}

func (a *Adapter) GetByKind(resources []manifest.ParsedResource, kind string) []manifest.ParsedResource {
	var result []manifest.ParsedResource
	for _, res := range resources {
		if res.Kind() == kind {
			result = append(result, res)
		}
	}
	return result
}

func (a *Adapter) GetSteps(resources []manifest.ParsedResource) []*manifest.Step {
	var result []*manifest.Step
	for _, res := range resources {
		if res.Step != nil {
			result = append(result, res.Step)
		}
	}
	return result
}

func (a *Adapter) GetCmds(resources []manifest.ParsedResource) []*manifest.Cmd {
	var result []*manifest.Cmd
	for _, res := range resources {
		if res.Cmd != nil {
			result = append(result, res.Cmd)
		}
	}
	return result
}

func (a *Adapter) GetCmdWorkflows(resources []manifest.ParsedResource) []*manifest.CmdWorkflow {
	var result []*manifest.CmdWorkflow
	for _, res := range resources {
		if res.CmdWorkflow != nil {
			result = append(result, res.CmdWorkflow)
		}
	}
	return result
}

func (a *Adapter) IndexByKind(resources []manifest.ParsedResource) (
	steps map[string]*manifest.Step,
	cmds map[string]*manifest.Cmd,
	cmdWorkflows map[string]*manifest.CmdWorkflow,
) {
	steps = make(map[string]*manifest.Step)
	cmds = make(map[string]*manifest.Cmd)
	cmdWorkflows = make(map[string]*manifest.CmdWorkflow)

	for _, res := range resources {
		switch {
		case res.Step != nil:
			steps[res.Step.Metadata.Name] = res.Step
		case res.Cmd != nil:
			cmds[res.Cmd.Metadata.Name] = res.Cmd
		case res.CmdWorkflow != nil:
			cmdWorkflows[res.CmdWorkflow.Metadata.Name] = res.CmdWorkflow
		}
	}

	return steps, cmds, cmdWorkflows
}

type BuildResult struct {
	Resources []manifest.ParsedResource
}

func NewBuildResult(resources []manifest.ParsedResource) *BuildResult {
	result := &BuildResult{
		Resources: []manifest.ParsedResource{},
	}

	for _, res := range resources {
		result.Resources = append(result.Resources, res)
	}

	logger.Info("kustomize", fmt.Sprintf("%d resources parsed", len(result.Resources)))

	return result
}
