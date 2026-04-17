package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type ProjectConfig struct {
	Name      string                 `json:"name"`
	Template  string                 `json:"template"`
	Variables map[string]interface{} `json:"variables"`
	Rewrites  map[string]string      `json:"rewrites"`
	Parse     []string               `json:"parse"`
}

type MenuEntry struct {
	Title  string `json:"title"`
	Href   string `json:"href"`
	Weight int    `json:"weight"`
	Sort   string `json:"sort"`
	Alias  string `json:"alias,omitempty"`
}

type Project struct {
	Config     ProjectConfig
	Name       string
	Path       string
	Template   string
	Variables  map[string]interface{}
	Rewrites   map[string]string
	ParseFiles []string
	Menus      map[string][]MenuEntry
	OutputPath string
}

func NewProject(path string) *Project {
	configPath := filepath.Join(path, "config.json")
	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Der Projekt-Ordner \"%s\" enthält keine config.json Datei\n", path)
		return nil
	}

	var config ProjectConfig
	if err := json.Unmarshal(contents, &config); err != nil {
		fmt.Printf("Der Projekt-Ordner \"%s\" enthält eine ungültige config.json Datei\n", path)
		return nil
	}

	return &Project{
		Config:     config,
		Name:       config.Name,
		Path:       path,
		Template:   config.Template,
		Variables:  config.Variables,
		Rewrites:   config.Rewrites,
		Menus:      map[string][]MenuEntry{"main": {}},
		OutputPath: filepath.Join("projects", config.Name, "output"),
	}
}

func (p *Project) GetMenus() {
	for _, file := range p.ParseFiles {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		html := string(content)

		// Title
		title := "! [h1] fehlt"
		reTitle := regexp.MustCompile("(?i)<title>(.*?)</title>")
		if match := reTitle.FindStringSubmatch(html); len(match) > 1 {
			title = match[1]
		}

		// Meta weight
		weight := 0
		reWeight := regexp.MustCompile(`(?i)data-weight="(\d+)"`)
		if match := reWeight.FindStringSubmatch(html); len(match) > 1 {
			fmt.Sscanf(match[1], "%d", &weight)
		}

		// Meta sort
		sortVal := ""
		reSort := regexp.MustCompile(`(?i)data-sort="([^"]*)"`)
		if match := reSort.FindStringSubmatch(html); len(match) > 1 {
			sortVal = match[1]
		}

		href := file[len(p.OutputPath):]
		// Clean up href separator (cross-platform fix)
		href = strings.TrimPrefix(href, string(os.PathSeparator))
		href = filepath.ToSlash(href)

		entry := MenuEntry{
			Title:  title,
			Href:   href,
			Weight: weight,
			Sort:   sortVal,
		}

		// Meta menu
		reMenu := regexp.MustCompile(`(?i)data-menu="([^"]*)"`)
		menuMatches := reMenu.FindAllStringSubmatch(html, -1)
		if len(menuMatches) == 0 {
			p.Menus["main"] = append(p.Menus["main"], entry)
		} else {
			for _, m := range menuMatches {
				menuName := m[1]
				p.Menus[menuName] = append(p.Menus[menuName], entry)
			}
		}

		// Meta alias
		reAlias := regexp.MustCompile(`(?i)data-alias="([^"]*)"`)
		aliasMatches := reAlias.FindAllStringSubmatch(html, -1)
		for _, am := range aliasMatches {
			aliasEntry := entry
			aliasEntry.Alias = am[1]
			p.Menus["main"] = append(p.Menus["main"], aliasEntry)
		}
	}
}

func (p *Project) RenderMenu(name, className, current string) string {
	entries, ok := p.Menus[name]
	if !ok {
		return ""
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Weight < entries[j].Weight
	})

	var out strings.Builder
	out.WriteString(fmt.Sprintf("<ul class=\"%s\">", className))
	for _, entry := range entries {
		active := ""
		if entry.Href == current {
			active = " class=\"active\" "
		}
		out.WriteString(fmt.Sprintf("<li><a %shref=\"%s\">%s</a></li>", active, entry.Href, entry.Title))
	}
	out.WriteString("</ul>")
	return out.String()
}

func (p *Project) RenderMenuTopic(name, className, current string) string {
	var topics []string
	if topicsVar, ok := p.Variables["topics"].([]interface{}); ok {
		for _, t := range topicsVar {
			topics = append(topics, fmt.Sprint(t))
		}
	}
	topics = append(topics, "???")

	allEntries := p.Menus[name]
	for i := range allEntries {
		if allEntries[i].Sort == "" {
			allEntries[i].Sort = "???"
		}
	}

	var out strings.Builder
	for _, topic := range topics {
		var topicEntries []MenuEntry
		for _, e := range allEntries {
			if e.Sort == topic {
				topicEntries = append(topicEntries, e)
			}
		}

		if topic == "???" && len(topicEntries) == 0 {
			continue
		}

		sort.Slice(topicEntries, func(i, j int) bool {
			return topicEntries[i].Weight < topicEntries[j].Weight
		})

		out.WriteString(fmt.Sprintf("<h3>%s</h3><ul>", topic))
		for _, e := range topicEntries {
			active := ""
			if e.Href == current {
				active = " class=\"active\" "
			}
			out.WriteString(fmt.Sprintf("<li><a %shref=\"%s\">%s</a></li>", active, e.Href, e.Title))
		}
		out.WriteString("</ul>")
	}

	return out.String()
}

func (p *Project) RenderMenuAlphabetical(name, className, current string) string {
	entries := make([]MenuEntry, len(p.Menus[name]))
	copy(entries, p.Menus[name])

	for i := range entries {
		title := entries[i].Title
		if entries[i].Alias != "" {
			title = entries[i].Alias
		}
		entries[i].Sort = sanitizeForSort(title)
	}

	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Sort) < strings.ToLower(entries[j].Sort)
	})

	var out strings.Builder
	out.WriteString("<ul>")
	letter := ""
	for _, item := range entries {
		if len(item.Sort) > 0 {
			firstChar := strings.ToUpper(string(item.Sort[0]))
			if letter != firstChar {
				if letter != "" {
					out.WriteString("</ul>")
				}
				letter = firstChar
				out.WriteString(fmt.Sprintf("<h2>%s</h2>", letter))
				out.WriteString("<ul>")
			}
		}

		active := ""
		if item.Href == current {
			active = " class=\"active\" "
		}
		aliasText := ""
		if item.Alias != "" {
			aliasText = fmt.Sprintf("<span style=\"color:gray\">%s ⟼ </span>", item.Alias)
		}
		out.WriteString(fmt.Sprintf("<li>%s<a %shref=\"%s\">%s</a></li>", aliasText, active, item.Href, item.Title))
	}
	out.WriteString("</ul>")
	return out.String()
}

func sanitizeForSort(s string) string {
	s = strings.ReplaceAll(s, "Ä", "Ae")
	s = strings.ReplaceAll(s, "ä", "ae")
	s = strings.ReplaceAll(s, "Ö", "Oe")
	s = strings.ReplaceAll(s, "ö", "oe")
	s = strings.ReplaceAll(s, "Ü", "Ue")
	s = strings.ReplaceAll(s, "ü", "ue")
	s = strings.ReplaceAll(s, "²", "2")
	reg, _ := regexp.Compile("[^0-9A-Za-z]")
	return reg.ReplaceAllString(s, "")
}

func (p *Project) CreateSearch() {
	fmt.Println("Erstelle Suche...")
	searchData := make(map[string]map[string]interface{})

	for _, file := range p.ParseFiles {
		fmt.Println(file)
		content, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		html := string(content)

		filename := file[len(p.OutputPath):]
		filename = strings.TrimPrefix(filename, string(os.PathSeparator))
		filename = filepath.ToSlash(filename)

		title := ""
		reTitle := regexp.MustCompile("(?i)<title>(.*?)</title>")
		if match := reTitle.FindStringSubmatch(html); len(match) > 1 {
			title = match[1]
		} else {
			Printh(fmt.Sprintf("Es scheint als hätte die Datei \"%s\" kein <title>-Tag.", file))
		}

		var tags []string
		reTags := regexp.MustCompile(`(?i)<meta\s+name="keywords"\s+content="([^"]*)"`)
		if match := reTags.FindStringSubmatch(html); len(match) > 1 {
			tags = strings.Split(match[1], ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		var aliases []string
		reAlias := regexp.MustCompile(`(?i)data-alias="([^"]*)"`)
		matches := reAlias.FindAllStringSubmatch(html, -1)
		for _, m := range matches {
			aliases = append(aliases, m[1])
		}

		searchData[filename] = map[string]interface{}{
			"title":   title,
			"tags":    tags,
			"aliases": aliases,
		}
	}

	jsonData, _ := json.Marshal(searchData)
	ioutil.WriteFile(filepath.Join(p.OutputPath, "search.json"), jsonData, 0644)
}

func (p *Project) Render() string {
	fmt.Println("Deleting old files...")
	os.RemoveAll(p.OutputPath)

	fmt.Println("Copy source files...")
	CopyDir(filepath.Join(p.Path, "src"), p.OutputPath)

	// Get all files to be parsed based on glob patterns
	p.ParseFiles = nil
	for _, pattern := range p.Config.Parse {
		matches, _ := filepath.Glob(filepath.Join(p.OutputPath, pattern))
		for _, match := range matches {
			if !strings.HasPrefix(filepath.Base(match), "._") {
				p.ParseFiles = append(p.ParseFiles, match)
			}
		}
	}

	template := NewTemplate(p.Template, p.Variables)
	template.Copy(p.OutputPath)

	p.GetMenus()
	p.CreateSearch()

	for _, file := range p.ParseFiles {
		Clear()
		fmt.Println("Erstelle Datei:")
		fmt.Println(file)

		currentHref := file[len(p.OutputPath):]
		currentHref = strings.TrimPrefix(currentHref, string(os.PathSeparator))
		currentHref = filepath.ToSlash(currentHref)

		if _, hasTopics := p.Variables["topics"]; !hasTopics {
			template.SetVariable("mainmenu", p.RenderMenuAlphabetical("main", "", currentHref))
		} else {
			template.SetVariable("mainmenu", p.RenderMenuTopic("main", "", currentHref))
		}
		template.SetVariable("topmenu", p.RenderMenu("top", "topmenu", currentHref))

		parsed, err := template.Parse(file)
		if err == nil {
			ioutil.WriteFile(file, []byte(parsed), 0644)
		}
	}

	return p.OutputPath
}
