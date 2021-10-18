package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	PrintPositive("Welcome to Banquet")
	flag.Parse()
	arguments := flag.Args()
	argumentCount := len(arguments)

	if argumentCount < 1 {
		printHelp("")
	}

	banquetOperation := arguments[0]

	switch banquetOperation {
		case "init":
			if argumentCount > 1 {
				initOperation := arguments[1]
				if initOperation == "kitchen" {

				}
			} else {
				for {
					banquetLocation := UserInput("Where would you like to serve banquet applications? (firebase, aws, localhost): ")
					if banquetLocation == "firebase" {
						fmt.Println("Checking for serviceAccountKey.json...")
						break
					} else if banquetLocation == "aws" {
						break
					} else if banquetLocation == "localhost" {
						break
					}
				}

			}
		case "dish":
			if argumentCount < 2 {
				printHelp("dish")
			}
			dishOperation := arguments[1]
			if argumentCount < 3 {
				printHelp("dish")
			} else {
				dishID := arguments[2]
				if dishOperation == "add" {
					dishTitle := UserInput("Please enter a title for your application: ")

					dishURL := UserInput("Please enter a GitHub URL for banquet to locate your application: ")

					var imageURLs []string

					for {
						imageURL := UserInput("Please enter an image URL followed by -n to add more images: ")
						if strings.Contains(imageURL, " -n") {
							imageURL = strings.TrimSuffix(imageURL, " -n")
							imageURLs = append(imageURLs, imageURL)
						} else {
							imageURLs = append(imageURLs, imageURL)
							break
						}
					}

					var colors []string

					for {
						color := UserInput("Please enter a color hex code followed by -n to add more colors: ")
						if strings.Contains(color, " -n") {
							color = strings.TrimSuffix(color, " -n")
							colors = append(colors, color)
						} else {
							colors = append(colors, color)
							break
						}
					}

					dish := dish{
						ID: dishID,
						Title: dishTitle,
						URL: dishURL,
						ImageURLs: imageURLs,
						Colors: colors,
						Status: "stopped",
					}
					addDish(dish)

				}
				if dishOperation == "get" {
					fmt.Println(getDish(dishID))
				}
				if dishOperation == "remove" {
					removeDish(dishID)
				}
				if dishOperation == "serve" {
					serveDish(dishID)
				}
			}
			if dishOperation == "get" {
				fmt.Println(getDishes())
			}

		default:
			printHelp("")
	}
}


// Prints out a list of viable commands and information about the project
func printHelp(command string) {
	switch command {
		case "dish":
			PrintNegative("To use the dish command ...")
		case "init":
			PrintNegative("To use the init command ...")
		default:
			PrintPositive("Available commands: ")
	}
	os.Exit(0)
}