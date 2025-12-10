package main

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
)

type Data struct {
	ProjectName string
}

func runCmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func writeTemplate(fileSystem embed.FS, name, destination string, data Data) error {
	bytes, err := fileSystem.ReadFile(filepath.Join("templates", name))
	if err != nil {
		return err
	}
	tmpl, err := template.New(name).Parse(string(bytes))
	if err != nil {
		return err
	}
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	return tmpl.Execute(outputFile, data)
}

//go:embed resources/*
var resourcesFS embed.FS

//go:embed templates/*
var templatesFS embed.FS

func printHelp() {
	fmt.Println("Initialize a new project: gyatt init <project-name>")
}

func initProject() {
	if len(os.Args) != 3 {
		printHelp()
		return
	}
	projectName := os.Args[2]
	data := Data{
		ProjectName: projectName,
	}
	fmt.Println("Initializing git repository")
	if err := runCmd("git", "init"); err != nil {
		fmt.Println("Failed to initialize git repository: ", err.Error())
		return
	}
	if err := runCmd("go", "mod", "init", projectName); err != nil {
		fmt.Println("Failed to init Go module: ", err.Error())
		return
	}
	//Create the .gitignore
	if err := writeTemplate(
		templatesFS, "gitignore.gotmpl",
		".gitignore", data,
	); err != nil {
		fmt.Println("Failed to create gitignore: ", err.Error())
	}
	//Create the directories
	for _, dir := range []string{"handler", "db", "ui"} {
		if err := os.Mkdir(dir, 0755); err != nil {
			fmt.Println("Failed to create folder: ", err.Error())
			return
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "init":
		initProject()
	default:
		printHelp()
	}
}
