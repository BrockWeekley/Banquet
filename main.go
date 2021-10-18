package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type user struct {
	GithubUsername string
	DeploymentType string
	ServiceAccountKey string
}

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
				// https://github.com/{OWNER}/{REPOSITORY}/archive/{BRANCH_NAME}.zip
				gitUser := UserInput("Please provide your GitHub Username: ")
				var banquetLocation string
				serviceAccountKeyLocation := ""
				for {
					banquetLocation = UserInput("Where would you like to serve Banquet applications? (firebase, aws, localhost): ")
					if banquetLocation == "firebase" {
						fmt.Println("In order to use Firebase with Banquet, you must generate a serviceAccountKey.json in the Firebase Console.")
						fmt.Println("https://firebase.google.com/docs/admin/setup#initialize-sdk")
						serviceAccountKeyLocation = UserInput("Please provide the path to your serviceAccountKey.json on this machine: ")
						break
					} else if banquetLocation == "aws" {
						fmt.Println("In order to use AWS with Banquet, you must generate an AWS Access Key in the AWS Management Console.")
						fmt.Println("https://aws.github.io/aws-sdk-go-v2/docs/getting-started/#get-your-aws-access-keys")
						serviceAccountKeyLocation = UserInput("Please provide the path to your new_user_credentials.csv on this machine: ")
						break
					} else if banquetLocation == "localhost" {
						break
					}
				}

				file, err := os.ReadFile("./config.json")
				CheckForError(err)
				var user user
				CheckForError(json.Unmarshal(file, &user))
				user.GithubUsername = gitUser
				user.DeploymentType = banquetLocation
				user.ServiceAccountKey = serviceAccountKeyLocation
				userBytes, err := json.Marshal(user)
				CheckForError(err)
				err = os.WriteFile("./config.json", userBytes, 0666)
			}
		case "dish":
			if argumentCount < 2 {
				printHelp("dish")
			}

			file, err := os.ReadFile("./config.json")
			if err != nil {
				printHelp("init")
			}
			var user user
			CheckForError(json.Unmarshal(file, &user))

			dishOperation := arguments[1]
			if argumentCount < 3 {
				if dishOperation == "get" {
					foundDishes := getDishes()
					for _, currentDish := range foundDishes {
						fmt.Println("ID: " + currentDish.ID + ", Title: " + currentDish.Title + ", Deployment Type: " + currentDish.DeploymentType + ", Status: " + currentDish.Status)
					}
				} else {
					printHelp("dish")
				}
			} else {
				dishID := arguments[2]
				if dishOperation == "add" {
					dishTitle := UserInput("Please enter a title for your application: ")

					dishRepository := UserInput("Please enter a GitHub Repository name for banquet to locate your application: ")

					dishBranch := UserInput("Please enter a GitHub Branch for banquet to locate your application (blank for 'master'): ")

					var dishToken string
					for {
						privateStatus := UserInput("Is your repository private or public? (private, public): ")
						if privateStatus == "private" {
							fmt.Println("You will need to generate a GitHub Personal Access Token to allow Banquet to access your repository.")
							dishToken = UserInput("Please provide your GitHub Personal Access Token: ")
							break
						}
						if privateStatus == "public" {
							dishToken = ""
							break
						}
					}

					if dishBranch == "" {
						dishBranch = "master"
					}

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

					url := ""
					if dishToken != "" {
						url = `https://github.com/` + user.GithubUsername + `/` + dishRepository + `/archive/` + dishBranch + `.zip` + `?` + dishToken
					} else {
						url = `https://github.com/` + user.GithubUsername + `/` + dishRepository + `/archive/` + dishBranch + `.zip`
					}

					localhostName := ""
					if user.DeploymentType == "firebase" {

					}
					if user.DeploymentType == "aws" {

					}
					if user.DeploymentType == "localhost" {
						localhostName = UserInput("Please provide the domain name that Banquet will route requests to: ")
					}

					dish := dish{
						ID: dishID,
						Title: dishTitle,
						URL: url,
						ImageURLs: imageURLs,
						Colors: colors,
						Status: "stopped",
						DeploymentType: user.DeploymentType,
						LocalhostName: localhostName,
					}
					addDish(dish)

				}
				if dishOperation == "get" {
					fmt.Println(getDish(dishID))
				}
				if dishOperation == "remove" {
					removeDish(dishID)
				}
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
			PrintNegative("You must run 'banquet init' before using Banquet. To use the init command type 'banquet init' in the console.")
			PrintNegative("You can also run 'banquet init kitchen' to start the kitchen API.")
		default:
			PrintPositive("Available commands: ")
	}
	os.Exit(0)
}