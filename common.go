package main

import (
	"fmt"
	"log"
)

// Color variables for displaying error messages

const Blue = "\033[1;34m%s\033[0m"
const Red = "\033[0;31m%s\033[0m"


/* Global Helper Functions */

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
	fmt.Printf(Blue, message + "\n")
}

func PrintNegative(message string) {
	fmt.Printf(Red, message + "\n")
}