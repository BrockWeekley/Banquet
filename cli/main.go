package main

import (
	"flag"
	common "github.com/brockweekley/banquet"
	kitchen "github.com/brockweekley/banquet/api"
	"github.com/go-git/go-git/v5"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Determine the passed in arguments
func main() {
	serveOption := flag.Bool("serve", false, "A flag used to skip the build process")
	flag.Parse()
	arguments := flag.Args()
	argumentCount := len(arguments)
	banquetOperation := arguments[0]

	if argumentCount < 2 {
		printHelp()
	}

	switch banquetOperation {

	case "course":

		common.PrintPositive("Add a new course:")
		if argumentCount > 2 {
			courseOperation := arguments[1]
			projectName := arguments[2]

			if courseOperation == "add" && projectName != "" {
				if argumentCount > 4 {
					githubURL := arguments[3]
					changeDirectory(projectName)
					directory, _ := os.Getwd()
					cloneRepository(githubURL, directory)
					common.PrintPositive("Course cloned and added to menu")
				}

			} else if courseOperation == "remove" && projectName != "" {
				changeDirectory(projectName)
			}
		}
		//												0		1		2				3
		common.PrintNegative("Invalid format: banquet course <option> <project_name> <repository_link?>")
		common.PrintNegative("Type 'banquet help' for more information.")

	case "reserve":
		port := arguments[1]
		common.PrintPositive("Reserving")
		if !*serveOption {
			install()
			build()
		}
		serve(port)

	default:
		printHelp()
	}
}

// Resolves the repository to the current project on any system, then creates a new folder for the project name
func changeDirectory(projectName string) {
	_, fileName, _, _ := runtime.Caller(0)
	filePath := strings.Trim(fileName, "/cli/main.go") + "/menu" + "/" + projectName
	common.CheckForError(os.MkdirAll(filePath, os.ModePerm))
	common.CheckForError(os.Chdir(filePath))
	_, pathRetrievalError := os.Getwd()
	common.CheckForError(pathRetrievalError)
}

// Clones the provides repository link to the newly created project folder
func cloneRepository(githubURL string, directory string) {
	_, cloneError := git.PlainClone(directory, false, &git.CloneOptions{
		URL: githubURL,
		Progress: os.Stdout,
	})
	common.CheckForError(cloneError)
}

// Installs the node modules for the waiter app
func install() {
	command := exec.Command("npm", "install")
	command.Dir = "./web/waiter"
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	common.PrintPositive("Running npm install for waiter, outputting stack:")
	common.CheckForError(command.Run())
	common.PrintPositive("Install ran correctly")
}

// Builds the waiter app
func build() {
	command := exec.Command("npm", "run", "build")
	command.Dir = "./web/waiter"
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	common.PrintPositive("Running npm build for waiter, outputting stack:")
	common.CheckForError(command.Run())
	common.PrintPositive("Build ran correctly")
}

// Serves the waiter app and kitchen API
func serve(port string) {
	fileServer := http.FileServer(http.Dir("./web/waiter/build/"))
	http.Handle("/", fileServer)
	go func() {
		kitchen.StartServer()
	}()
	common.PrintPositive("Waiter is serving...")
	common.CheckForError(http.ListenAndServe(":" + port, nil))
}

// Prints out a list of viable commands and information about the project
func printHelp() {
	common.PrintPositive("Help: ")
}