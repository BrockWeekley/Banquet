package main

import (
	. "github.com/brockweekley/banquet"
	"github.com/go-git/go-git"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Determine the passed in arguments
func main() {
	arguments := os.Args[1:]
	argumentCount := len(arguments)
	banquetOperation := arguments[0]

	if argumentCount < 1 {
		printHelp()
	}

	switch banquetOperation {

	case "course":

		PrintPositive("Add a new course:")

		courseOperation := arguments[1]

		if argumentCount > 3 && courseOperation != "" && arguments[2] != "" {

			if courseOperation == "add" {
				githubURL := arguments[2]
				projectName := arguments[3]

				if argumentCount > 4 && projectName != "" {
					changeDirectory(projectName)
					directory, _ := os.Getwd()
					cloneRepository(githubURL, directory)
					PrintPositive("Course cloned and added to menu")
				}

			} else if courseOperation == "remove" {
				projectName := arguments[2]

				changeDirectory(projectName)
			}

		}

		PrintNegative("Invalid format: banquet course <option> <repository_link?> <project_name>")

	case "reserve":
		PrintPositive("Reserving")
		install()
		build()
		serve()

	default:
		printHelp()
	}
}

// Resolves the repository to the current project on any system, then creates a new folder for the project name
func changeDirectory(projectName string) {
	_, fileName, _, _ := runtime.Caller(0)
	filePath := strings.Trim(fileName, "/cli/main.go") + "/menu" + "/" + projectName
	CheckForError(os.MkdirAll(filePath, os.ModePerm))
	CheckForError(os.Chdir(filePath))
	_, pathRetrievalError := os.Getwd()
	CheckForError(pathRetrievalError)
}

// Clones the provides repository link to the newly created project folder
func cloneRepository(githubURL string, directory string) {
	_, cloneError := git.PlainClone(directory, false, &git.CloneOptions{
		URL: githubURL,
		Progress: os.Stdout,
	})
	CheckForError(cloneError)
}

// Installs the node modules for the waiter app
func install() {
	command := exec.Command("npm", "install")
	command.Dir = "./web/waiter/"
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	PrintPositive("Running npm install for waiter, outputting stack:")
	CheckForError(command.Run())
	PrintPositive("Install ran correctly")
}

// Builds the waiter app
func build() {
	command := exec.Command("npm", "run", "build")
	command.Dir = "./web/waiter"
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	PrintPositive("Running npm build for waiter, outputting stack:")
	CheckForError(command.Run())
	PrintPositive("Build ran correctly")
}

// Serves the waiter app
func serve() {
	fileServer := http.FileServer(http.Dir("./web/waiter/build/"))
	http.Handle("/", fileServer)
	PrintPositive("Waiter is serving...")
	CheckForError(http.ListenAndServe(":8080", nil))
}

// Prints out a list of viable commands and information about the project
func printHelp() {
	PrintPositive("Help: ")
}