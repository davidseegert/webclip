package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	ProjectsPath string `json:"projectsPath"`
}

type ProjectManager struct {
	Projects      map[string]*Project
	ActiveProject *Project
}

func NewProjectManager() *ProjectManager {
	return &ProjectManager{
		Projects: make(map[string]*Project),
	}
}

func (pm *ProjectManager) Init() error {
	// Loading configuration
	contents, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("=== ERROR ===")
		fmt.Println("The file \"config.json\" was not found in the main directory.")
		fmt.Println("1) Make sure the file exists.")
		fmt.Println("2) Make sure the script was started from the main directory.")
		fmt.Println("==============")
		return err
	}

	var config Config
	err = json.Unmarshal(contents, &config)
	if err != nil {
		fmt.Println("=== ERROR ===")
		fmt.Println("The file \"config.json\" in the main directory could not be read correctly.")
		fmt.Println("Make sure the file contains only valid JSON.")
		fmt.Println("==============")
		return err
	}

	// Scan projects directory
	entries, err := ioutil.ReadDir(config.ProjectsPath)
	if err != nil {
		// If the projects path doesn't exist, we just have no projects
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			projectPath := filepath.Join(config.ProjectsPath, entry.Name())
			if _, err := os.Stat(filepath.Join(projectPath, "config.json")); err == nil {
				pm.AddProject(projectPath)
			}
		}
	}

	return nil
}

func (pm *ProjectManager) AddProject(path string) {
	proj := NewProject(path)
	if proj != nil {
		pm.Projects[proj.Name] = proj
	}
}

func (pm *ProjectManager) GetProject(name string) *Project {
	if proj, ok := pm.Projects[name]; ok {
		pm.ActiveProject = proj
		return proj
	}
	return nil
}

func (pm *ProjectManager) ProjectSelector() *Project {
	for name := range pm.Projects {
		fmt.Println(name)
	}

	for {
		fmt.Print("\nEnter name of project to compile: ")
		var projectName string
		fmt.Scanln(&projectName)

		if proj := pm.GetProject(projectName); proj != nil {
			return proj
		}
		fmt.Printf("\"%s\" is no project.\n", projectName)
	}
}
