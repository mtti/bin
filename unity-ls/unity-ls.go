/*
Copyright 2021 Matti Hiltunen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
	"gopkg.in/yaml.v2"
)

type projectVersionFile struct {
	EditorVersion string `yaml:"m_EditorVersion"`
}

type project struct {
	Path    string
	Version *version.Version
}

func readVersion(directory string) (*version.Version, error) {
	bytes, err := ioutil.ReadFile(path.Join(directory, "ProjectVersion.txt"))
	if err != nil {
		return nil, err
	}

	parsed := projectVersionFile{}
	err = yaml.Unmarshal(bytes, &parsed)
	if err != nil {
		return nil, err
	}

	v, err := version.NewVersion(parsed.EditorVersion)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func scanDirectory(directory string, recursive bool, level int) []project {
	projects := make([]project, 0)

	v, err := readVersion(path.Join(directory, "ProjectSettings"))
	if err == nil {
		projects = append(projects, project{
			Path:    directory,
			Version: v,
		})
		return projects
	}

	if !recursive && level > 0 {
		return projects
	}

	children, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	subdirectories := make([]string, 0)

	for _, child := range children {
		name := child.Name()

		if child.IsDir() {
			if strings.HasPrefix(name, ".") {
				continue
			}

			subdirectories = append(subdirectories, name)
		}
	}

	for _, subdirectory := range subdirectories {
		subProjects := scanDirectory(path.Join(directory, subdirectory), recursive, level+1)
		projects = append(projects, subProjects...)
	}

	return projects
}

func main() {
	projects := scanDirectory(".", false, 0)
	byVersion := make(map[string][]string)
	versions := make([]*version.Version, 0)

	for _, p := range projects {
		if _, ok := byVersion[p.Version.Original()]; !ok {
			byVersion[p.Version.Original()] = make([]string, 0)
			versions = append(versions, p.Version)
		}
		byVersion[p.Version.Original()] = append(byVersion[p.Version.Original()], p.Path)
	}

	sort.Sort(version.Collection(versions))

	for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
		versions[i], versions[j] = versions[j], versions[i]
	}

	for _, v := range versions {
		fmt.Printf("%s\n", v.String())

		for _, projectPath := range byVersion[v.Original()] {
			fmt.Printf("    %s\n", projectPath)
		}

		fmt.Printf("\n")
	}
}
