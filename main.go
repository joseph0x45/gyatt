package main

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
)

type Data struct {
	ProjectName string
}

type projectFile struct {
	Name             string
	FileSystem       embed.FS
	WriteDestination string
}

func runCmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func writeTemplate(fileSystem embed.FS, name, destination string, data *Data, parseTemplate bool) error {
	bytes, err := fileSystem.ReadFile(name)
	if err != nil {
		return err
	}
	if parseTemplate {
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
	input, _ := fileSystem.ReadFile(name)
	return os.WriteFile(destination, input, 0755)
}

//go:embed resources/*
var resourcesFS embed.FS

//go:embed templates/*
var templatesFS embed.FS

func printHelp() {
	fmt.Println("Initialize a new project: gyatt init <project-name>")
	fmt.Println("Add a new dependency to the project: gyatt add-dependency <dependency-name>")
	fmt.Println("Available dependencies are: 'htmx', 'alpine', 'toastify'")
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
	//Create the directories
	for _, dir := range []string{"handler", "db", "ui/layouts", "static"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("Failed to create folder: ", err.Error())
			return
		}
	}
	projectFiles := []projectFile{
		{
			Name:             "templates/main.gotmpl",
			WriteDestination: "main.go",
			FileSystem:       templatesFS,
		},
		{
			Name:             "templates/db.gotmpl",
			WriteDestination: "db/db.go",
			FileSystem:       templatesFS,
		},
		{
			Name:             "templates/input.css",
			WriteDestination: "static/input.css",
			FileSystem:       templatesFS,
		},
		{
			Name:             "templates/gitignore.gotmpl",
			WriteDestination: ".gitignore",
			FileSystem:       templatesFS,
		},
		{
			Name:             "templates/Makefile",
			WriteDestination: "Makefile",
			FileSystem:       templatesFS,
		},
		{
			Name:             "templates/main_layout.gotmpl",
			WriteDestination: "ui/layouts/main.templ",
			FileSystem:       templatesFS,
		},
		{
			Name:             "templates/index.gotmpl",
			WriteDestination: "ui/index.templ",
			FileSystem:       templatesFS,
		},
		{
			Name:             "templates/go.sum",
			WriteDestination: "go.sum",
			FileSystem:       templatesFS,
		},
	}
	for _, file := range projectFiles {
		if err := writeTemplate(
			file.FileSystem, file.Name,
			file.WriteDestination, &data,
			true,
		); err != nil {
			fmt.Printf("Failed to create %s: %s\n", file.Name, err.Error())
			break
		}
	}
	if err := runCmd("go", "mod", "tidy"); err != nil {
		fmt.Println("Failed to run 'go mod tidy': ", err.Error())
		return
	}
	if err := runCmd("make", "build"); err != nil {
		fmt.Println("Failed to build project: ", err.Error())
		return
	}
}

func addDependency() {
	if len(os.Args) != 3 {
		printHelp()
		return
	}
	dependency := os.Args[2]
	switch dependency {
	case "htmx", "alpine":
		if err := writeTemplate(
			resourcesFS,
			fmt.Sprintf("resources/%s.js", dependency),
			fmt.Sprintf("static/%s.js", dependency),
			nil,
			false,
		); err != nil {
			fmt.Printf("Failed to add dependency %s: %s\n", dependency, err.Error())
			return
		}
	case "toastify":
		for _, dep := range []string{"toastify.css", "toastify.js"} {
			if err := writeTemplate(
				resourcesFS,
				fmt.Sprintf("resources/%s", dep),
				fmt.Sprintf("static/%s", dep),
				nil,
				false,
			); err != nil {
				fmt.Printf("Failed to add dependency %s: %s\n", dependency, err.Error())
				return
			}
		}
	default:
		printHelp()
	}
}

func setup() {
	if os.Geteuid() == 0 {
		data, err := resourcesFS.ReadFile("resources/templ")
		os.WriteFile("/usr/local/bin/templ", data, 0755)
		if err != nil {
			fmt.Println("Failed to setup templ: ", err.Error())
			return
		}
		data, err = resourcesFS.ReadFile("resources/tailwindcss")
		os.WriteFile("/usr/local/bin/tailwindcss", data, 0755)
		if err != nil {
			fmt.Println("Failed to setup tailwindcss: ", err.Error())
			return
		}
		return
	}
	fmt.Println("Run this command as root")
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "setup":
		setup()
	case "init":
		initProject()
	case "add-dependency":
		addDependency()
	default:
		printHelp()
	}
}
