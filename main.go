package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type user struct {
	GithubUsername string
	DeploymentType string
	ServiceAccountKey string
	GithubOAuthKey string
}

const banquetID = "b58241f56afaa752c830"

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
					wg := new(sync.WaitGroup)
					wg.Add(1)

					go func() {
						StartServer()
						wg.Done()
					}()

					wg.Wait()
				}
			} else {
				PrintPositive("Let's walk you through your first setup of Banquet\n")
				PrintPositive("Please ensure node package manager is installed on this machine before continuing the setup.")
				PrintPositive("Banquet uses Docker to create a reusable image of your application. In order to use Banquet locally, you will need to install Docker on this machine before adding any applications.")
				fmt.Println("https://docs.docker.com/get-docker/")
				//PrintPositive("\nIf you plan to use Banquet with a Google Cloud or Firebase account, the Cloud SDK will need to be installed on this machine before adding any applications.")
				//fmt.Println("https://cloud.google.com/sdk/docs/install#linux")
				PrintPositive("If you plan to use Banquet with an AWS account, the AWS CLI will need to be installed on this machine before adding any applications.")
				fmt.Println("https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html")
				PrintPositive("\nNow, please provide some information to get started. You can change this information in the future by rerunning the init command, or manually changing the config.json file.\n")

				gitUser := UserInput("Please provide your GitHub Username: ")
				var banquetLocation string
				serviceAccountKeyLocation := ""
				for {
					banquetLocation = UserInput("Where would you like to serve Banquet applications? (aws, localhost): ")
					//if banquetLocation == "gcloud" || banquetLocation == "firebase" {
					//	fmt.Println("In order to use Google Cloud or Firebase with Banquet, you must generate a serviceAccountKey.json and enable REST APIs in the Firebase or Google Cloud Console.")
					//	fmt.Println("https://firebase.google.com/docs/admin/setup#initialize-sdk")
					//	fmt.Println("https://firebase.google.com/docs/hosting/api-deploy#enable-api")
					//	serviceAccountKeyLocation = UserInput("Please provide the path to your serviceAccountKey.json on this machine: ")
					//	break
					// } else
					if banquetLocation == "aws" {
						fmt.Println("In order to use AWS with Banquet, you must generate an AWS Access Key in the AWS Management Console.")
						fmt.Println("https://aws.github.io/aws-sdk-go-v2/docs/getting-started/#get-your-aws-access-keys")
						serviceAccountKeyLocation = UserInput("Please provide the path to your new_user_credentials.csv on this machine: ")
						break
					} else if banquetLocation == "localhost" {
						break
					}
				}
				UpdateUser(gitUser, banquetLocation, serviceAccountKeyLocation, "", true)
				PrintPositive("User config has been updated. Installing required packages, this may take a while.")
				cmd := exec.Command("npm", "install", "typescript", "-g")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				CheckForError(cmd.Run())
				cmd = exec.Command("npm", "install", "react-scripts@latest", "-g")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				CheckForError(cmd.Run())
				PrintPositive("Happy dining!")
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
					foundDishes := GetDishes()
					for _, currentDish := range foundDishes {
						fmt.Println("ID: " + currentDish.ID + ", Title: " + currentDish.Title + ", Deployment Type: " + currentDish.DeploymentType + ", Status: " + currentDish.Status)
					}
				} else {
					printHelp("dish")
				}
			} else {
				dishID := arguments[2]
				if dishOperation == "add" {
					fmt.Println("Please ensure Docker is installed on the local machine before adding a dish.")
					if CheckForExistingDishID(dishID) {
						PrintNegative("A dish with this ID already exists. Run 'banquet dish remove {dishID}' to remove it.")
						return
					}

					existingBuild := UserInput("If you have an existing build you would like to use, please enter the title of that application (blank for no preexisting build):")


					dishTitle := UserInput("Please enter a title for this application: ")
					var sb strings.Builder
					_, err := sb.WriteString(dishTitle + " " + strings.ToLower(dishID))
					CheckForError(err)
					dishTitle = sb.String()

					var dishRepository string
					var dishToken string

					if existingBuild == "" {

						dishRepository = UserInput("Please enter the GitHub Repository name for banquet to locate your application (must be exactly as it appears on GitHub): ")

						//dishBranch := UserInput("Please enter a GitHub Branch for banquet to locate your application (blank for 'master'): ")

						for {
							privateStatus := UserInput("Is your repository private or public? (private, public): ")
							if privateStatus == "private" {
								if user.GithubOAuthKey != "" {
									fmt.Println("Using existing account key...")
									dishToken = user.GithubOAuthKey
								} else {
									fmt.Println("You will need to allow Banquet access to your repository:")
									data := map[string]string{"client_id": banquetID, "scope": "repo"}
									jsonData, err := json.Marshal(data)
									CheckForError(err)
									response, err := http.Post("https://github.com/login/device/code", "application/json", bytes.NewBuffer(jsonData))
									CheckForError(err)

									body, err := ioutil.ReadAll(response.Body)
									CheckForError(err)
									values := strings.Split(string(body), "&")
									deviceCode := ""
									userCode := ""
									verURL := ""

									for index, value := range values {
										if index == 0 {
											deviceCode = strings.Split(value, "=")[1]
										}
										if index == 3 {
											userCode = strings.Split(value, "=")[1]
										}
										if index == 4 {
											verURL = strings.Split(value, "=")[1]
										}
									}
									decodedURL, err := url.QueryUnescape(verURL)
									fmt.Println("Please navigate to: " + decodedURL + " and enter the following code: ")
									fmt.Println(userCode)

									UserInput("Press Enter when you have successfully authenticated.")

									defer CheckForError(response.Body.Close())

									data = map[string]string{"client_id": banquetID, "device_code": deviceCode, "grant_type": "urn:ietf:params:oauth:grant-type:device_code"}
									jsonData, err = json.Marshal(data)
									CheckForError(err)
									response, err = http.Post("https://github.com/login/oauth/access_token", "application/json", bytes.NewBuffer(jsonData))
									CheckForError(err)

									body, err = ioutil.ReadAll(response.Body)
									CheckForError(err)
									params := strings.Split(string(body), "&")
									dishToken = strings.Split(params[0], "=")[1]
									PrintPositive("You have successfully authenticated with GitHub.")

									save := UserInput("Would you like to save this token for all apps on this user account?")
									save = strings.ToLower(save)
									if save == "yes" || save == "y" || save == "ye" || save == "yeah" || save == "-y" {
										UpdateUser("", "", "", dishToken, false)
									}

									defer CheckForError(response.Body.Close())

									CheckForError(err)
								}
								break
							}
							if privateStatus == "public" {
								dishToken = ""
								break
							}
						}
					}


					ionicVariables := [9]string{"", "", "", "", "", "", "", "", ""}
					ionic := UserInput("Are you using Ionic themes for your project?")
					if ionic == "yes" || ionic == "y" || ionic == "ye" || ionic == "yeah" || ionic == "-y" {
						ionicVariables[0] = UserInput("Please enter a hex value for ion-color-primary. Example: '#ffffff' (Leave blank for no change):")
						ionicVariables[1] = UserInput("Please enter a hex value for ion-color-secondary. Example: '#ffffff' (Leave blank for no change):")
						ionicVariables[2] = UserInput("Please enter a hex value for ion-color-tertiary. Example: '#ffffff' (Leave blank for no change):")
						ionicVariables[3] = UserInput("Please enter a hex value for ion-color-success. Example: '#ffffff' (Leave blank for no change):")
						ionicVariables[4] = UserInput("Please enter a hex value for ion-color-warning. Example: '#ffffff' (Leave blank for no change):")
						ionicVariables[5] = UserInput("Please enter a hex value for ion-color-danger. Example: '#ffffff' (Leave blank for no change):")
						ionicVariables[6] = UserInput("Please enter a hex value for ion-color-dark. Example: '#ffffff' (Leave blank for no change):")
						ionicVariables[7] = UserInput("Please enter a hex value for ion-color-medium. Example: '#ffffff' (Leave blank for no change):")
						ionicVariables[8] = UserInput("Please enter a hex value for ion-color-light. Example: '#ffffff' (Leave blank for no change):")
					}

					capacitorStatus := ""
					capacitorResponse := UserInput("Would you like to use Capacitor to build your application for Android?")
					if capacitorResponse == "yes" || capacitorResponse == "y" || capacitorResponse == "ye" || capacitorResponse == "yeah" || capacitorResponse == "-y" {
						capacitorStatus = UserInput("Please provide the path to your SDK location on this machine (Ex: C\\:\\\\Users\\\\User\\\\AppData\\\\Local\\\\Android\\\\Sdk")
					}


					var customStyleLocation string
					custom := UserInput("Would you like to implement a custom stylesheet?")
					if custom == "yes" || custom == "y" || custom == "ye" || custom == "yeah" || custom == "-y" {
						customStyleLocation = UserInput("Please enter the location of your custom stylesheet:")
						break
					}

					var customTSLocation string
					custom = UserInput("Would you like to implement a custom typescript file?")
					if custom == "yes" || custom == "y" || custom == "ye" || custom == "yeah" || custom == "-y" {
						customTSLocation = UserInput("Please enter the location of your custom typescript file:")
						break
					}

					var imageURLs []string
					for {
						imageURL := UserInput("Please enter an image URL. To add more images after this one, type -n after the URL (Leave blank for no additions): ")
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
						color := UserInput("Please enter a color hex code or rgb function. Examples: '#ffffff' or 'rgb(255, 255, 255)' To add more colors after this one, type -n after the color (Leave blank for no additions): ")
						if strings.Contains(color, " -n") {
							color = strings.TrimSuffix(color, " -n")
							colors = append(colors, color)
						} else {
							colors = append(colors, color)
							break
						}
					}

					localhostName := ""
					if user.DeploymentType == "firebase" {
						// TODO: TBD
					}
					if user.DeploymentType == "aws" {
						// TODO: TBD
					}
					if user.DeploymentType == "localhost" {
						localhostName = UserInput("Please provide the port that banquet should deploy the container to: ")
					}

					dish := Dish{
						ID: dishID,
						ContainerID: "",
						Title: dishTitle,
						URL: `https://api.github.com/repos/` + user.GithubUsername + `/` + dishRepository + `/zipball/master`,
						ImageURLs: imageURLs,
						Colors: colors,
						CustomStyleLocation: customStyleLocation,
						CustomTSLocation: customTSLocation,
						IonicVariables: ionicVariables,
						Capacitor: capacitorStatus,
						Status: "stopped",
						DeploymentType: user.DeploymentType,
						LocalhostName: localhostName,
						Token: dishToken,
					}
					AddDish(dish, existingBuild)
				}
				if dishOperation == "get" {
					fmt.Println(GetDish(dishID))
				}
				if dishOperation == "remove" {
					RemoveDish(dishID)
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

func UpdateUser(gitUser string, banquetLocation string, serviceAccountKeyLocation string, authKey string, create bool) {
	file, err := os.ReadFile("./config.json")
	CheckForError(err)
	var user user
	CheckForError(json.Unmarshal(file, &user))
	if create {
		user.GithubUsername = gitUser
		user.DeploymentType = banquetLocation
		user.ServiceAccountKey = serviceAccountKeyLocation
		user.GithubOAuthKey = authKey
	} else {
		if gitUser != "" {
			user.GithubUsername = gitUser
		}
		if banquetLocation != "" {
			user.DeploymentType = banquetLocation
		}
		if serviceAccountKeyLocation != "" {
			user.ServiceAccountKey = serviceAccountKeyLocation
		}
		if authKey != "" {
			user.GithubOAuthKey = authKey
		}
	}

	userBytes, err := json.Marshal(user)
	CheckForError(err)
	err = os.WriteFile("./config.json", userBytes, 0666)
}