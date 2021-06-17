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

var Blue = "\033[1;34m%s\033[0m"
var Red = "\033[0;31m%s\033[0m"

func main() {
	arguments := os.Args[1:]
	if len(arguments) < 1 {
		arguments = append(arguments, "help")
	}
	switch arguments[0] {
	case "course":
		fmt.Printf(Blue, "Add a new course: \n")
		if len(arguments) > 3 && arguments[1] != "" && arguments[2] != "" && arguments[3] != "" {
			if arguments[1] == "add" {
				directoryError := changeDirectory(arguments[3])
				if directoryError == nil {
					directory, _ := os.Getwd()
					cloneError := cloneRepository(arguments[2], directory)
					if cloneError == nil {
						fmt.Printf(Blue, "Course cloned and added to menu")
					} else {
						log.Fatal(cloneError)
					}
				} else {
					log.Fatal(directoryError)
				}
			} else if arguments[1] == "remove" {

			}

		}

		fmt.Printf(Red, "Invalid format: banquet course <option> <repository_link> <project_name>")

	case "reserve":
		fmt.Println("Reserving")
		install()
		build()
		serve()
	default:
		fmt.Println("Help: ")
	}
}

func changeDirectory(projectName string) error {
	_, fileName, _, _ := runtime.Caller(0)
	filePath := strings.Trim(fileName, "/cli/main.go") + "/menu" + "/" + projectName
	folderCreationError := os.MkdirAll(filePath, os.ModePerm)
	directoryError := os.Chdir(filePath)
	_, pathRetrievalError := os.Getwd()
	if folderCreationError != nil {
		log.Fatal(folderCreationError)
	}
	if pathRetrievalError != nil {
		log.Fatal(pathRetrievalError)
	}
	return directoryError
}

func cloneRepository(githubURL string, directory string) error {
	_, cloneError := git.PlainClone(directory, false, &git.CloneOptions{
		URL: githubURL,
		Progress: os.Stdout,
	})
	return cloneError
}

func install() {
	command := exec.Command("npm", "install")
	command.Dir = "./web/waiter/"
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	fmt.Printf(Blue, "Running npm install for waiter, outputting stack: \n")
	commandError := command.Run()
	if commandError == nil {
		fmt.Printf(Blue, "Install ran correctly \n")
	} else {
		log.Fatal(commandError)
		return
	}
}

func build() {
	command := exec.Command("npm", "run", "build")
	command.Dir = "./web/waiter"
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	fmt.Printf(Blue, "Running npm build for waiter, outputting stack: \n")
	commandError := command.Run()
	if commandError == nil {
		fmt.Printf(Blue, "Build ran correctly \n")
	} else {
		log.Fatal(commandError)
		return
	}
}

func serve() {
	fileServer := http.FileServer(http.Dir("./web/waiter/build/"))
	http.Handle("/", fileServer)
	fmt.Printf(Blue, "Waiter is serving...")
	httpError := http.ListenAndServe(":8080", nil)
	if httpError != nil {
		log.Fatal(httpError)
		return
	}
}