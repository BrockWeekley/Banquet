package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var Blue = "\033[1;34m%s\033[0m"

func main() {
	install()
	build()
	serve()
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