package main

import (
	"fmt"
	"github.com/go-git/go-git"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Color variables for displaying error messages

const Blue = "\033[1;34m%s\033[0m"
const Red = "\033[0;31m%s\033[0m"

// Determine the passed in arguments
func main() {
	arguments := os.Args[1:]
	if len(arguments) < 1 {
		printHelp()
	}

	switch arguments[0] {

	case "course":
		PrintPositive("Add a new course:")
		if len(arguments) > 3 && arguments[1] != "" && arguments[2] != "" && arguments[3] != "" {
			if arguments[1] == "add" {
				changeDirectory(arguments[3])
				directory, _ := os.Getwd()
				cloneRepository(arguments[2], directory)
				PrintPositive("Course cloned and added to menu")
			} else if arguments[1] == "remove" {
				changeDirectory(arguments[3])
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


/* Global Helper Functions */

func CheckForError(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
}

func PrintPositive(message string) {
	fmt.Printf(Blue, message + "\n")
}

func PrintNegative(message string) {
	fmt.Printf(Red, message + "\n")
}