package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
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

		PrintPositive("Add a new course:")
		if argumentCount > 2 {
			courseOperation := arguments[1]
			projectName := arguments[2]

			if courseOperation == "add" && projectName != "" {
				if argumentCount > 3 {
					githubURL := arguments[3]
					changeDirectory(projectName)
					directory, _ := os.Getwd()
					cloneRepository(githubURL, directory)
					initializeStatus(projectName)
					PrintPositive("Course cloned and added to menu")
					break
				}

			} else if courseOperation == "remove" && projectName != "" {
				changeDirectory(projectName)
				break
			}
		}
		//												0		1		2				3
		PrintNegative("Invalid format: banquet course <option> <project_name> <repository_link?>")
		PrintNegative("Type 'banquet help' for more information.")

	case "reserve":
		port := arguments[1]
		PrintPositive("Reserving")
		if !*serveOption {
			install()
			build()
		}
		serve(port)

	case "serve":
		PrintPositive("Serve a previously added course here")
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

func initializeStatus(projectName string) {
	_, fileName, _, _ := runtime.Caller(0)
	filePath := strings.ReplaceAll(fileName, "/cli/main.go", "") + "/api/"
	CheckForError(os.Chdir(filePath))
	menu, fileError := os.ReadFile("menu.json")
	CheckForError(fileError)
	var courses []Course
	CheckForError(json.Unmarshal(menu, &courses))
	var newMenu []Course
	newMenu = append(courses, Course{Name: projectName, Status: "Stopped", Port: 0})
	initialize, jsonError := json.Marshal(newMenu)
	CheckForError(jsonError)
	CheckForError(os.WriteFile("menu.json", initialize, 0666))
}

// Installs the node modules for the waiter app
func install() {
	command := exec.Command("npm", "install")
	command.Dir = "./web/waiter"
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

// Serves the waiter app and kitchen API
func serve(port string) {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	serveMux := http.NewServeMux()
	serveServer := http.Server {
		Addr: fmt.Sprintf(":%v", port),
		Handler: serveMux,
	}
	fileServer := http.FileServer(http.Dir("./web/waiter/build/"))
	serveMux.Handle("/", fileServer)
	go func() {
		PrintPositive("Waiter is serving...")
		CheckForError(serveServer.ListenAndServe())
		wg.Done()
	}()

	go func() {
		StartServer()
		wg.Done()
	}()

	wg.Wait()

}
//func(w http.ResponseWriter, r *http.Request) {
//	path, err := os.Getwd()
//	if err != nil {
//		log.Println(err)
//	}
//	suffix := string(path[len(path)-3:])
//	fmt.Println(suffix)
//	if suffix == "api" {
//		http.StripPrefix("/api", fileServer).ServeHTTP(w, r)
//	} else {
//		fileServer.ServeHTTP(w, r)
//	}
//}


// Prints out a list of viable commands and information about the project
func printHelp() {
	PrintPositive("Help: ")
}