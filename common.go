package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
)

// Color variables for displaying error messages

const Blue = "\033[1;34m%s\033[0m"
const Red = "\033[0;31m%s\033[0m"


/* Global Helper Functions */

func UserInput(prompt string)(response string) {
	fmt.Println(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	CheckForError(err)
	response = strings.TrimRight(response, "\r\n")
	response = strings.TrimRight(response, "\r")
	response = strings.TrimRight(response, "\n")
	response = strings.ToLower(response)
	response = strings.TrimSpace(response)

	if strings.Compare(response, "cancel") == 0 {
		fmt.Println("Cancelled")
		os.Exit(0)
	}
	return response
}

func CheckForError(err error) {
	if err == nil {
	return
	}
	log.Fatal(err)
}

func ThrowError(message string) {
	log.Fatal(message)
}

func PrintPositive(message string) {
	color.Blue(message + "\n")
}

func PrintNegative(message string) {
	color.Red(message + "\n")
}