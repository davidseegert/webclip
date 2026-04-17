package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type TemplateConfig struct {
	Parse     []string          `json:"parse"`
	Copy      []string          `json:"copy"`
	Variables map[string]interface{} `json:"variables"`
}

type Template struct {
	Name      string
	Variables map[string]interface{}
	Config    TemplateConfig
}

func NewTemplate(name string, variables map[string]interface{}) *Template {
	t := &Template{
		Name:      name,
		Variables: make(map[string]interface{}),
	}

	templatePath := filepath.Join("templates", name)
	configPath := filepath.Join(templatePath, "config.json")
	
	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error loading template config for %s: %v\n", name, err)
		return nil
	}

	var config TemplateConfig
	if err := json.Unmarshal(contents, &config); err != nil {
		fmt.Printf("Error parsing template config for %s: %v\n", name, err)
		return nil
	}

	t.Config = config
	for k, v := range config.Variables {
		t.Variables[k] = v
	}
	if variables != nil {
		for k, v := range variables {
			t.Variables[k] = v
		}
	}

	return t
}

func (t *Template) SetVariable(varName string, value interface{}) {
	t.Variables[varName] = value
}

func (t *Template) Copy(destFolder string) error {
	// Copy folders from template
	fmt.Println("Copy folders from template...")
	for _, folder := range t.Config.Copy {
		src := filepath.Join("templates", t.Name, folder)
		dst := filepath.Join(destFolder, folder)
		if err := CopyDir(src, dst); err != nil {
			return err
		}
	}

	// Parsing files from template
	fmt.Println("Begin parsing template files...")
	for _, folder := range t.Config.Parse {
		src := filepath.Join("templates", t.Name, folder)
		dst := filepath.Join(destFolder, folder)
		os.MkdirAll(dst, 0755)

		files, _ := ioutil.ReadDir(src)
		for _, f := range files {
			if f.IsDir() || strings.HasPrefix(f.Name(), "._") {
				continue
			}
			srcFile := filepath.Join(src, f.Name())
			dstFile := filepath.Join(dst, f.Name())

			fmt.Printf("Parsing template file: %s\n", dstFile)
			content, err := ioutil.ReadFile(srcFile)
			if err != nil {
				continue
			}
			rendered := t.RenderString(string(content), t.Variables)
			ioutil.WriteFile(dstFile, []byte(rendered), 0644)
		}
	}
	return nil
}

func (t *Template) Parse(filePath string) (string, error) {
	data, err := t.GetData(filePath)
	if err != nil {
		return "", err
	}

	templateFile := filepath.Join("templates", t.Name, "index.html")
	content, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return "", err
	}

	return t.RenderString(string(content), data), nil
}

// GetData mimics Template.rb's getData: extracts title, body, and data- attributes
func (t *Template) GetData(filePath string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	for k, v := range t.Variables {
		data[k] = v
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	html := string(content)

	// Extract title
	reTitle := regexp.MustCompile("(?i)<title>(.*?)</title>")
	if match := reTitle.FindStringSubmatch(html); len(match) > 1 {
		data["title"] = match[1]
	}

	// Extract body
	reBody := regexp.MustCompile("(?is)<body>(.*?)</body>")
	if match := reBody.FindStringSubmatch(html); len(match) > 1 {
		data["body"] = match[1]
	}

	// Extract data- attributes from meta tags
	reMeta := regexp.MustCompile(`(?i)<meta\s+[^>]*>`)
	metaTags := reMeta.FindAllString(html, -1)
	
	reDataAttr := regexp.MustCompile(`(?i)data-([^=]+)="([^"]*)"`)

	for _, meta := range metaTags {
		matches := reDataAttr.FindAllStringSubmatch(meta, -1)
		for _, m := range matches {
			key := m[1]
			value := m[2]
			var val interface{} = value
			if strings.Contains(value, ";") {
				val = strings.Split(value, ";")
			}
			data[key] = val
		}
	}

	return data, nil
}

// RenderString is the "mini-Liquid" engine
func (t *Template) RenderString(content string, variables map[string]interface{}) string {
	result := content
	// We handle {{ key }} and {{key}}
	re := regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)
	
	// Direct replacement for simplicity (Standard Library only)
	return re.ReplaceAllStringFunc(result, func(match string) string {
		key := re.FindStringSubmatch(match)[1]
		if val, ok := variables[key]; ok {
			switch v := val.(type) {
			case string:
				return v
			case []string:
				return strings.Join(v, ", ")
			default:
				return fmt.Sprintf("%v", v)
			}
		}
		return match // Keep original if not found
	})
}
