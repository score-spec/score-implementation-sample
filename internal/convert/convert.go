// Copyright 2024 Humanitec
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package convert

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/score-spec/score-go/framework"
	scoretypes "github.com/score-spec/score-go/types"
	"gopkg.in/yaml.v3"

	"github.com/score-spec/score-implementation-sample/internal/state"
)

func Workload(currentState *state.State, workloadName string) (map[string]interface{}, error) {
	resOutputs, err := currentState.GetResourceOutputForWorkload(workloadName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate outputs: %w", err)
	}
	sf := framework.BuildSubstitutionFunction(currentState.Workloads[workloadName].Spec.Metadata, resOutputs)

	spec := currentState.Workloads[workloadName].Spec
	containers := maps.Clone(spec.Containers)
	for containerName, container := range containers {
		if container.Variables, err = convertContainerVariables(container.Variables, sf); err != nil {
			return nil, fmt.Errorf("workload: %s: container: %s: variables: %w", workloadName, containerName, err)
		}

		if container.Files, err = convertContainerFiles(container.Files, currentState.Workloads[workloadName].File, sf); err != nil {
			return nil, fmt.Errorf("workload: %s: container: %s: files: %w", workloadName, containerName, err)
		}
		containers[containerName] = container
	}
	spec.Containers = containers
	resources := maps.Clone(spec.Resources)
	for resName, res := range resources {
		resUid := framework.NewResourceUid(workloadName, resName, res.Type, res.Class, res.Id)
		resState, ok := currentState.Resources[resUid]
		if !ok {
			return nil, fmt.Errorf("workload '%s': resource '%s' (%s) is not primed", workloadName, resName, resUid)
		}
		res.Params = resState.Params
		resources[resName] = res
	}
	spec.Resources = resources

	// ===============================================================================
	// TODO: HERE IS WHERE YOU MAY CONVERT THE WORKLOAD INTO YOUR TARGET MANIFEST TYPE
	// ===============================================================================

	raw, err := yaml.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("workload: %s: failed to serialise manifest: %w", workloadName, err)
	}
	var intermediate map[string]interface{}
	_ = yaml.Unmarshal(raw, &intermediate)

	return intermediate, nil
}

func convertContainerVariables(input scoretypes.ContainerVariables, sf func(string) (string, error)) (map[string]string, error) {
	outMap := make(map[string]string, len(input))
	for key, value := range input {
		out, err := framework.SubstituteString(value, sf)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}
		outMap[key] = out
	}
	return outMap, nil
}

func convertContainerFiles(input []scoretypes.ContainerFilesElem, scoreFile *string, sf func(string) (string, error)) ([]scoretypes.ContainerFilesElem, error) {
	outSlice := make([]scoretypes.ContainerFilesElem, 0, len(input))
	for i, fileElem := range input {
		var content string
		if fileElem.Content != nil {
			content = *fileElem.Content
		} else if fileElem.Source != nil {
			sourcePath := *fileElem.Source
			if !filepath.IsAbs(sourcePath) && scoreFile != nil {
				sourcePath = filepath.Join(filepath.Dir(*scoreFile), sourcePath)
			}
			if rawContent, err := os.ReadFile(sourcePath); err != nil {
				return nil, fmt.Errorf("%d: source: failed to read file '%s': %w", i, sourcePath, err)
			} else {
				content = string(rawContent)
			}
		} else {
			return nil, fmt.Errorf("%d: missing 'content' or 'source'", i)
		}

		var err error
		if fileElem.NoExpand == nil || !*fileElem.NoExpand {
			content, err = framework.SubstituteString(string(content), sf)
			if err != nil {
				return nil, fmt.Errorf("%d: failed to substitute in content: %w", i, err)
			}
		}
		fileElem.Source = nil
		fileElem.Content = &content
		bTrue := true
		fileElem.NoExpand = &bTrue
		outSlice = append(outSlice, fileElem)
	}
	return outSlice, nil
}
