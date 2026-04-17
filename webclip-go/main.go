package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	Clear()

	pm := NewProjectManager()
	if err := pm.Init(); err != nil {
		// Error messages are already printed in Init
		os.Exit(1)
	}

	// Manually select project
	project := pm.ProjectSelector()
	if project == nil {
		fmt.Println("No project selected.")
		os.Exit(1)
	}

	outpath := project.Render()

	Clear()
	fmt.Println("+-------------------------------+")
	fmt.Println("| Projekt erfolgreich erstellt! |")
	fmt.Println("+-------------------------------+")
	fmt.Println("")
	fmt.Println("Output-Ordner:")
	
	absPath, _ := filepath.Abs(".")
	fmt.Printf("%s/%s\n", absPath, outpath)
	fmt.Println("")
}
