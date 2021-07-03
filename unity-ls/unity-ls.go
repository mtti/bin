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

func readVersion(directory string) (string, error) {
	bytes, err := ioutil.ReadFile(path.Join(directory, "ProjectVersion.txt"))
	if err != nil {
		return "", err
	}

	parsed := projectVersionFile{}
	err = yaml.Unmarshal(bytes, &parsed)
	if err != nil {
		panic(err)
	}

	return parsed.EditorVersion, nil
}

func scanDirectory(directory string) []project {
	projects := make([]project, 0)

	children, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	subdirectories := make([]string, 0)

	for _, child := range children {
		name := child.Name()
		childPath := path.Join(directory, name)

		if child.IsDir() {
			if strings.HasPrefix(name, ".") {
				continue
			}

			if name == "ProjectSettings" {
				rawVersion, err := readVersion(childPath)

				if err != nil {
					fmt.Printf("[WARN] %s: %s\n", childPath, err)
					continue
				}

				v, err := version.NewVersion(rawVersion)
				if err != nil {
					fmt.Printf("[WARN] %s: %s\n", childPath, err)
					continue
				}

				projects = append(projects, project{
					Path:    directory,
					Version: v,
				})

				continue
			}
			subdirectories = append(subdirectories, name)
		}
	}

	for _, subdirectory := range subdirectories {
		subProjects := scanDirectory(path.Join(directory, subdirectory))
		projects = append(projects, subProjects...)
	}

	return projects
}

func main() {
	projects := scanDirectory(".")
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
