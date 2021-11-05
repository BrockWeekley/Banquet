package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type dish struct {
	ID string
	Title string
	URL string
	ImageURLs []string
	Colors []string
	Status string
	DeploymentType string
	LocalhostName string
	Token string
}

func getDishes()(dishes []dish) {
	file, err := os.ReadFile("./menu.json")
	CheckForError(err)
	CheckForError(json.Unmarshal(file, &dishes))

	return dishes
}

func getDish(dishID string)(foundDish dish) {
	file, err := os.ReadFile("./menu.json")
	CheckForError(err)
	var dishes []dish
	CheckForError(json.Unmarshal(file, &dishes))
	for _, currentDish := range dishes {
		if currentDish.ID == dishID {
			foundDish = currentDish
		}
	}

	return foundDish
}

func checkForExistingDishID(potentialDishID string)(status bool) {
	file, err := os.ReadFile("./menu.json")
	var dishes []dish
	CheckForError(json.Unmarshal(file, &dishes))
	CheckForError(err)
	for _, dish := range dishes {
		if dish.ID == potentialDishID {
			return true
		}
	}
	return false
}

func addDish(newDish dish)() {
	file, err := os.ReadFile("./menu.json")
	var dishes []dish
	CheckForError(json.Unmarshal(file, &dishes))
	CheckForError(err)
	newDish.Title = strings.ReplaceAll(newDish.Title, " ", "_")
	dishes = append(dishes, newDish)
	dishBytes, err := json.Marshal(dishes)
	CheckForError(err)
	err = os.WriteFile("./menu.json", dishBytes, 0644)
	CheckForError(err)
	serveDish(newDish)
}

func removeDish(dishID string)(status bool) {
	file, err := os.ReadFile("./menu.json")
	CheckForError(err)
	var dishes []dish
	var foundDish dish
	CheckForError(json.Unmarshal(file, &dishes))
	for i, currentDish := range dishes {
		if currentDish.ID == dishID {
			foundDish = currentDish
			dishes = append(dishes[:i], dishes[i+1:]...)
		}
	}
	dishBytes, err := json.Marshal(dishes)
	CheckForError(err)
	err = os.WriteFile("./menu.json", dishBytes, 0666)
	CheckForError(err)
	cleanDish(foundDish)
	return foundDish.ID == dishID
}

func serveDish(dish dish)() {
	file, err := os.ReadFile("./config.json")
	CheckForError(err)
	var user user
	CheckForError(json.Unmarshal(file, &user))

	downloadRepo(dish)
	generateStyling(dish)
	dockerize(dish)

	if dish.DeploymentType == "firebase" {

	}
	if dish.DeploymentType == "aws" {

	}
	if dish.DeploymentType == "localhost" {

	}
}

func cleanDish(dish dish)() {
	file, err := os.ReadFile("./config.json")
	CheckForError(err)
	var user user
	CheckForError(json.Unmarshal(file, &user))

	if dish.DeploymentType == "firebase" {

	}
	if dish.DeploymentType == "aws" {

	}
	if dish.DeploymentType == "localhost" {

	}
}

func downloadRepo(dish dish) {
	if _, err := os.Stat("./menu"); os.IsNotExist(err) {
		CheckForError(os.Mkdir("./menu", os.ModeDir))
	}
	create, err := os.Create("./menu/" + dish.Title + ".zip")
	CheckForError(err)

	if dish.Token != "" {
		client := http.Client{}
		request, err := http.NewRequest("GET", dish.URL, nil)
		CheckForError(err)
		request.Header.Add("Authorization", "token " + dish.Token)
		response, err := client.Do(request)
		CheckForError(err)
		if response.StatusCode != http.StatusOK {
			PrintNegative("Bad Status for provided GitHub URL: " + response.Status)
		}
		_, err = io.Copy(create, response.Body)

		defer CheckForError(response.Body.Close())
	} else {
		client := http.Client{}
		request, err := http.NewRequest("GET", dish.URL, nil)
		CheckForError(err)
		response, err := client.Do(request)
		CheckForError(err)
		if response.StatusCode != http.StatusOK {
			PrintNegative("Bad Status for provided GitHub URL: " + response.Status)
		}
		_, err = io.Copy(create, response.Body)

		defer CheckForError(response.Body.Close())
	}

	reader, err := zip.OpenReader("./menu/" + dish.Title + ".zip")
	destination := "./menu/" + dish.Title + "/"
	CheckForError(err)

	var filenames []string

	for _, foundFile := range reader.File {

		path := filepath.Join(destination, foundFile.Name)

		filenames = append(filenames, path)

		if foundFile.FileInfo().IsDir() {
			CheckForError(os.MkdirAll(path, os.ModePerm))
			continue
		}

		CheckForError(os.MkdirAll(filepath.Dir(path), os.ModePerm))

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, foundFile.Mode())
		CheckForError(err)

		rc, err := foundFile.Open()
		CheckForError(err)

		_, err = io.Copy(outFile, rc)

		CheckForError(err)

		CheckForError(outFile.Close())
		CheckForError(rc.Close())
	}

	folders, err := os.ReadDir("./menu/" + dish.Title)
	CheckForError(err)
	err = os.Rename("./menu/" + dish.Title + "/" + folders[0].Name(), "./menu/" + dish.Title + "/" + dish.Title)
	CheckForError(err)
	defer CheckForError(reader.Close())
	defer CheckForError(create.Close())
	defer CheckForError(os.Remove("./menu/" + dish.Title + ".zip"))
}

func generateStyling(dish dish) {
	css, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/src/banquet.css")
	CheckForError(err)
	ts, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/src/banquet.ts")
	CheckForError(err)
	for index, image := range dish.ImageURLs {
		if image != "" {
			_, err = css.WriteString(".img" + strconv.Itoa(index + 1) + "{background-image: url(" + image + ");}\n")
			CheckForError(err)
		}
	}
	for index, color := range dish.Colors {
		if color != "" {
			_, err = css.WriteString(".color" + strconv.Itoa(index + 1) + "{color: " + color + ";}\n")
			CheckForError(err)
			_, err = css.WriteString(".bcolor" + strconv.Itoa(index + 1) + "{background-color: " + color + ";}\n")
			CheckForError(err)
		}
	}
	_, err = ts.WriteString("export const title = \"Tech Co\";")
	defer CheckForError(css.Close())
	defer CheckForError(ts.Close())
}

func dockerize(dish dish) {
	dockerFile, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/Dockerfile")
	dockerIgnore, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/.dockerignore")
	CheckForError(err)
	_, err = dockerFile.WriteString("FROM node:16.13.0\n\n" +
								"WORKDIR /" + dish.Title + "\n\n" +
								"ENV PATH /" + dish.Title + "/node_modules/.bin:$PATH\n\n" +
								"COPY package.json ./\n" +
								"RUN npm install\n" +
								"RUN npm install react-scripts@latest -g\n\n" +
								"COPY . ./\n\n" +
								"RUN npm run build\n\n" +
								"CMD [\"npm\", \"start\"]")
	CheckForError(err)
	_, err = dockerIgnore.WriteString("node_modules\n" +
										"build\n" +
										".dockerignore\n" +
										"Dockerfile\n" +
										"Dockerfile.prod")
	CheckForError(err)

	//TODO: This is a hacky way to do this. We should be able to npm run build, take the relatively few files, tar them and use the Docker SDK
	_, fileName, _, _ := runtime.Caller(0)
	filePath := strings.ReplaceAll(fileName, "dishes.go", "")
	CheckForError(os.Chdir(filePath + "menu/" + dish.Title + "/" + dish.Title + "/"))

	PrintPositive("Running Docker command. This will take a few minutes...")
	cmd := exec.Command("docker", "build", "-t", dish.Title + ":latest", ".")
	CheckForError(err)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	CheckForError(cmd.Start())

	defer CheckForError(dockerFile.Close())
	defer CheckForError(dockerIgnore.Close())
}